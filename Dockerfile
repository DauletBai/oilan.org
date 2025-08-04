# oilan/Dockerfile

# Stage 1: Build the application
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/oilan ./cmd/server

# Stage 2: Create a minimal final image
FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/oilan .
# Copy the configs folder into our final image
COPY ./configs ./configs
COPY ./web ./web

EXPOSE 8080
CMD ["./oilan"]