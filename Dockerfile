# # Stage 1: Build
# FROM golang:1.25-alpine AS builder

# # Установим необходимые пакеты
# RUN apk add --no-cache git gcc musl-dev

# # Установим рабочую директорию
# WORKDIR /app

# # Копируем go.mod и go.sum
# COPY go.mod go.sum ./

# # Загружаем зависимости
# RUN go mod download

# # Копируем весь исходный код
# COPY . .

# # Собираем приложение
# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/main/main.go


# # Stage 2: Runtime
# FROM alpine:latest

# # Установим ca-certificates для HTTPS
# RUN apk --no-cache add ca-certificates
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# # Установим рабочую директорию
# WORKDIR /app

# # Копируем бинарник из builder stage
# COPY --from=builder /app/server .

# # Копируем .env файл если существует
# COPY .env* ./

# # Создаём entrypoint.sh
# RUN echo '#!/bin/sh\n\
# set -e\n\
# echo "Running migrations..."\n\
# goose -dir ./sql/migrations up\n\
# echo "Starting server..."\n\
# exec ./server' > ./entrypoint.sh && chmod +x /entrypoint.sh

# # Expose порт
# EXPOSE 8080

# # Запускаем приложение
# # CMD ["./server"]
# ENTRYPOINT ["./entrypoint.sh"]

# Stage 1: Build
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/main/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/sql ./sql

# Копируем .env файл если существует
COPY .env* ./

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]