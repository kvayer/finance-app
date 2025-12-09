# Builder stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/app/main.go

# Runner stage
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
# Копируем шаблоны и .env
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]