package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy/src/handlers"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "myao-proxy", log.LstdFlags)

	nh := handlers.NewNewsHandler(l)
	lh := handlers.NewLinksHandler(l)

	sm := mux.NewRouter()
	sm.HandleFunc("/news", nh.GetSimpleRequest)
	sm.HandleFunc("/tpu-news", nh.TPURequest)
	sm.HandleFunc("/get-links", lh.FindAllLinks)

	server := &http.Server{
		Addr:         ":9091",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9091")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown. Signal:", sig)

	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	server.Shutdown(tc)
}
