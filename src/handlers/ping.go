package handlers

import (
	"net/http"
	"proxy/src/loggers"
)

type Pinger struct {
	logger *loggers.AggregatedLoggers
}

func NewPinger(l *loggers.AggregatedLoggers) *Pinger {
	return &Pinger{l}
}

func (p *Pinger) Handle(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("[INFO] Ping request received")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from proxy!"))
}
