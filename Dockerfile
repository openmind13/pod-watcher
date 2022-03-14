FROM golang:1.17.8-alpine3.15 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN GOOS=linux GOARCH=amd64 go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o pod-watcher cmd/pod-watcher.go

FROM --platform=linux/amd64 alpine:3.15
WORKDIR /app
COPY --from=builder /build/pod-watcher .
CMD ["/app/pod-watcher"]