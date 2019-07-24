FROM golang:1.12.7-alpine as builder
WORKDIR /app
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lcp-go

FROM scratch
LABEL maintainer="levshino@gmail.com"
LABEL description="Simple and fast proxy to bypass CORS issues during prototyping against existing APIs"
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env /app/
COPY --from=builder /app/lcp-go /app/
EXPOSE 8118
ENTRYPOINT ["./lcp-go"]
