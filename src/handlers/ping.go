package handlers

import (
	"net/http"
	"proxy/src/loggers"
)

type Pinger struct {
	l *loggers.AggregatedLoggers
}

func NewPinger(l *loggers.AggregatedLoggers) *Pinger {
	return &Pinger{l}
}

func (p *Pinger) Ping(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from proxy!"))
}
