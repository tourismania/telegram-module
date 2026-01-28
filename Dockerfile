# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Установить необходимые инструменты
RUN apk add --no-cache git tzdata

WORKDIR /app

# Копировать go.mod и go.sum
COPY go.mod go.sum ./

# Загрузить зависимости
RUN go mod download

# Копировать исходный код
COPY . .

# Собрать приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

# Установить tzdata для временных зон
RUN apk --no-cache add tzdata

WORKDIR /app

# Копировать бинарный файл из builder
COPY --from=builder /app/telegram .

# Копировать миграции
# COPY migrations/ ./migrations/

# Создать непривилегированного пользователя
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
USER appuser

# Expose порт
EXPOSE 8080

# Health check
# HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
#    CMD /app/telegram health || exit 1

# Запустить приложение
CMD ["/app/telegram", "serve"]