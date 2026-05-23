FROM node:20-alpine AS frontend

WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN go build -o doraemon ./cmd/doraemon

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/doraemon .
RUN mkdir -p data/files

ENTRYPOINT ["./doraemon"]
