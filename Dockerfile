FROM golang:1.19

WORKDIR /usr/src/go-app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./
RUN go build -o /go-proxy

EXPOSE 9091

CMD ["/go-proxy"]