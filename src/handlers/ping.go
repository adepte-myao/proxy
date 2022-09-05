package handlers

import (
	"log"
	"net/http"
)

type Pinger struct {
	l *log.Logger
}

func NewPinger(l *log.Logger) *Pinger {
	return &Pinger{l}
}

func (p *Pinger) Ping(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from proxy!"))
}
