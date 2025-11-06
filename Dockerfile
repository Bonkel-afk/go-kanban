# 1. Build-Stage
FROM golang:1.22-alpine AS builder

# Arbeitsverzeichnis im Container
WORKDIR /app

# Go-Module vorbereiten
COPY go.mod go.sum ./
RUN go mod download

# Restlichen Code kopieren
COPY . .

# Statisches Linux-Binary bauen
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o kanban ./server

# 2. Runtime-Stage (kleines Image ohne Go-Toolchain)
FROM alpine:3.20

WORKDIR /app

# Binary aus dem Builder-Stage kopieren
COPY --from=builder /app/kanban ./kanban

# Standard-Port deines Servers
EXPOSE 9090

# Default: File-Storage (Mongo steuerst du über docker-compose-ENV)
ENV KANBAN_STORAGE=file

CMD ["./kanban"]
