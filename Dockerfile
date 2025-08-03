# oilan/Dockerfile

# Этап 1: Сборка приложения
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/oilan ./cmd/server

# Этап 2: Создание минималистичного образа
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/oilan .
# В будущем сюда можно будет скопировать папку с конфигами и фронтендом
# COPY --from=builder /app/configs ./configs
# COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./oilan"]