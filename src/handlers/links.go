package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"proxy/src/dto"
	"proxy/src/loggers"
)

type LinksHandler struct {
	logger *loggers.AggregatedLoggers
}

func NewLinksHandler(l *loggers.AggregatedLoggers) *LinksHandler {
	return &LinksHandler{l}
}

func (lh LinksHandler) FindAllLinks(rw http.ResponseWriter, r *http.Request) {
	lh.logger.Println("[INFO] Get links request")
	var rd dto.LinksRequestData
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		lh.logger.Println("[ERROR] Decoding failed, stop processing")

		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid object body, must be dto.LinksRequestData"))
		return
	}

	resp, err := http.Get(rd.Link)
	if err != nil {
		lh.logger.Println("[ERROR] Can't receive response from given source, stop processing", err)

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source wasn't received. Check your URL or try later"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lh.logger.Println("Status code is not OK, stop processing")

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is not OK"))
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		lh.logger.Println("[ERROR] Can't read response body")

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when reading response body"))
		return
	}

	bodyString := string(bodyBytes)
	entries := getAllHrefPartsFromStringifyBody(bodyString)

	rw.WriteHeader(http.StatusOK)
	for _, v := range entries {
		ref := getReferenceFromHref(v)
		rw.Write([]byte(ref))
		rw.Write([]byte("\n"))
	}
}

func getAllHrefPartsFromStringifyBody(str string) []string {
	reg := regexp.MustCompile(`href="[^"]*://[^"]*"`)
	return reg.FindAllString(str, -1)
}

func getReferenceFromHref(hrefMatch string) string {
	// hrefMatch is a string like `href="required-reference.org"`
	return hrefMatch[6 : len(hrefMatch)-1]
}
