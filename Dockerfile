FROM golang:1.19-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /bin/go-proxy

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build /bin/go-proxy /bin/go-proxy
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9091
ENTRYPOINT [ "/bin/go-proxy" ]