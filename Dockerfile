FROM golang:1.22-alpine

RUN apk update && apk upgrade

WORKDIR /app

EXPOSE 80
EXPOSE 443

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY cmd ./cmd
COPY templates ./templates

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/roxy.go

CMD ["./roxy"]