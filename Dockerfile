FROM golang:1.12.7-alpine as builder

LABEL maintainer="levshino@gmail.com"
LABEL description="Simple and fast proxy to bypass CORS issues during prototyping against existing APIs"

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lcp-go

EXPOSE 8118
ENTRYPOINT ["./lcp-go"]
