# --- Этап сборки ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем из папки cmd/chat
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o nofelet-chat ./cmd/chat

# --- Этап запуска ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем бинарь с новым именем
COPY --from=builder /app/nofelet-chat .

EXPOSE 8444

# Запускаем именно nofelet
CMD ["./nofelet-chat"]