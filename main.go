package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"proxy/src/handlers"
	"proxy/src/loggers"

	"github.com/gorilla/mux"
)

func main() {
	fo, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("Couldn't open the file, stop executing")
		return
	}

	filelog := log.New(fo, "myao-proxy", log.LstdFlags)
	conslog := log.New(os.Stdout, "myao-proxy", log.LstdFlags)
	logger := loggers.NewAggregatedLoggers(filelog, conslog)

	lh := handlers.NewLinksHandler(logger)
	ph := handlers.NewPinger(logger)

	sm := mux.NewRouter()
	sm.HandleFunc("/", ph.Handle)
	sm.HandleFunc("/get-links", lh.FindAllLinks)

	server := &http.Server{
		Addr:         ":9091",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	startInterruptedServerAsync(server, logger)
}

func startInterruptedServerAsync(server *http.Server, logger *loggers.AggregatedLoggers) {
	go func() {
		logger.Println("Starting server on port 9091")

		err := server.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	useUserInterrupt(server, logger)
}

func useUserInterrupt(server *http.Server, logger *loggers.AggregatedLoggers) {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown. Signal:", sig)

	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	server.Shutdown(tc)
}
