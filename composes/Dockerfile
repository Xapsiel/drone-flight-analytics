# Стадия сборки
FROM golang:1.24.1-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Диагностика: вывод списка файлов в контексте
RUN ls -la

# Копирование go.mod и go.sum (go.sum опционально)
COPY go.mod ./
COPY go.sum* ./
# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Финальная стадия
FROM alpine:3.18

# Установка необходимых пакетов
RUN apk add --no-cache ca-certificates tzdata postgis

# Создание непривилегированного пользователя
RUN adduser -D -g '' appuser

# Установка рабочей директории
WORKDIR /app

# Копирование бинарного файла
COPY --from=builder /app/server .

# Копирование конфигурации
COPY --from=builder /app/config ./config

# Копирование веб-ресурсов (если есть)

# Установка прав доступа
RUN chown -R appuser:appuser /app

# Переключение на непривилегированного пользователя
USER appuser

# Настройка окружения
ENV PORT=8080
ENV GIN_MODE=release

# Открытие порта
EXPOSE 8080

# Команда для запуска приложения
CMD ["./server", "-c", "/app/config/config.yaml"]