FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o doraemon ./cmd/doraemon

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/doraemon .
RUN mkdir -p data/files

ENTRYPOINT ["./doraemon"]
