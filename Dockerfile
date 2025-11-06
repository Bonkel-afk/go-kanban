# Multi-Stage Build: erst bauen, dann schlankes Image
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Go-Module abhängigkeiten laden
COPY go.mod ./
RUN go mod download

# Restlicher Code
COPY . .

# Binary bauen
RUN go build -o kanban .

# --- Runtime Image ---
FROM alpine:3.20

WORKDIR /app

# CA-Zertifikate (falls du später HTTP-Requests machst)
RUN apk add --no-cache ca-certificates

# Binary aus Builder kopieren
COPY --from=builder /app/kanban /app/kanban

# tasks.json liegt im Arbeitsverzeichnis /app
# (wird später per Volume gemountet, wenn du willst)
EXPOSE 9090

CMD ["/app/kanban"]
