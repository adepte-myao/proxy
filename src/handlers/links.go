package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"proxy/src/dto"
	"proxy/src/loggers"
	"regexp"
)

type LinksHandler struct {
	l *loggers.AggregatedLoggers
}

func NewLinksHandler(l *loggers.AggregatedLoggers) *LinksHandler {
	return &LinksHandler{l}
}

func (lh LinksHandler) FindAllLinks(rw http.ResponseWriter, r *http.Request) {
	lh.l.Println("[INFO] Get links request")

	var rd dto.LinksRequestData
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		lh.l.Println("[ERROR] Decoding failed, stop processing")

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid object body, must be dto.LinksRequestData"))
		return
	}

	resp, err := http.Get(rd.Link)
	if err != nil {
		lh.l.Println("[ERROR] Can't receive response from given source, stop processing")

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source wasn't received. Check your URL or try later"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lh.l.Println("Status code is not OK, stop processing")

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is not OK"))
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		lh.l.Println("[ERROR] Can't read response body")

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when reading response body"))
		return
	}

	bodyString := string(bodyBytes)

	reg := regexp.MustCompile(`href="[^"]*://[^"]*"`)
	entries := reg.FindAllString(bodyString, -1)

	rw.WriteHeader(http.StatusOK)
	for _, v := range entries {
		rw.Write([]byte(v[5:]))
		rw.Write([]byte("\n"))
	}
}
