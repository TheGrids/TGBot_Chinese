# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bot ./cmd/bot

# Stage 2: Minimal runtime
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bot .
COPY --from=builder /app/assets ./assets

# Запускаем бинарник
ENTRYPOINT ["./bot"]
