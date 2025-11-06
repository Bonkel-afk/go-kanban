FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o kanban ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/kanban /app/kanban
EXPOSE 9090
CMD ["/app/kanban"]
