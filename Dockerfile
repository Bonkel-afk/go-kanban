# ---------- Build Stage ----------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Go Modules
COPY go.mod go.sum ./
RUN go mod download

# Source Code
COPY . .

# Build Binary
RUN go build -o kanban ./server

# ---------- Runtime Stage ----------
FROM alpine:3.20

WORKDIR /app

# Binary kopieren
COPY --from=builder /app/kanban /app/kanban

# Templates & Static-Files mit ins Image
COPY --from=builder /app/internals/web/templates /app/internals/web/templates
COPY --from=builder /app/internals/web/static /app/internals/web/static

# Expose Port
EXPOSE 9090

# Default Environment (kann in docker-compose überschrieben werden)
ENV KANBAN_STORAGE=file

CMD ["/app/kanban"]
