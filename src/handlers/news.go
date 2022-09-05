package handlers

import (
	"io"
	"log"
	"net/http"
	"regexp"
)

type NewsHandler struct {
	l *log.Logger
}

func NewNewsHandler(l *log.Logger) *NewsHandler {
	return &NewsHandler{l}
}

func (h *NewsHandler) GetSimpleRequest(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("[INFO] Get request received")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from server\n"))
}

func (h *NewsHandler) TPURequest(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("[INFO] TPU Request")

	resp, err := http.Get("https://tpu.ru")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.l.Println("Status code is not OK, stop processing")
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.l.Println("Error when reading response body")
		return
	}

	bodyString := string(bodyBytes)

	reg := regexp.MustCompile("<a.*>.*</a>")
	entries := reg.FindAllString(bodyString, -1)

	rw.WriteHeader(http.StatusOK)
	for _, v := range entries {
		rw.Write([]byte(v))
		rw.Write([]byte("\n"))
	}
}
