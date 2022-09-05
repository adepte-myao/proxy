package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"proxy/src/dto"
	"regexp"
)

type LinksHandler struct {
	l *log.Logger
}

func NewLinksHandler(l *log.Logger) *LinksHandler {
	return &LinksHandler{l}
}

func (lh LinksHandler) FindAllLinks(rw http.ResponseWriter, r *http.Request) {
	lh.l.Println("[INFO] Get links request")

	var rd dto.LinksRequestData
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		lh.l.Println("[ERROR] Decoding failed, stop processing")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := http.Get(rd.Link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lh.l.Println("Status code is not OK, stop processing")
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		lh.l.Println("Error when reading response body")
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
