# Stage 1: Build the migrate CLI tool
FROM golang:1.25-alpine AS migrate_builder

# Install necessary dependencies for building CGO-enabled binaries (like for postgres)
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app

# Use Go modules to install the migrate CLI with appropriate database tags (e.g., 'postgres')
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ----------------------------------------------------------
# Stage 2: Build
FROM golang:1.25-alpine AS app_builder

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

# ----------------------------------------------------------
# Stage 3: Runtime
FROM alpine:latest

# Установить tzdata для временных зон
RUN apk --no-cache add tzdata

WORKDIR /app

# Copy the built migrate binary from the first stage
COPY --from=migrate_builder /go/bin/migrate /usr/local/bin/migrate
# Копировать бинарный файл из builder
COPY --from=app_builder /app/telegram .
# Копировать папку с миграциями
COPY database/postgres/migrations ./migrations

# Создать непривилегированного пользователя
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
USER appuser

# Expose порт
EXPOSE 8080

# Health check
# HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
#    CMD /app/telegram health || exit 1

COPY --from=app_builder /app/docker/docker-entrypoint.sh /app/docker-entrypoint

ENTRYPOINT ["./docker-entrypoint"]

# Запустить приложение
CMD ["/app/telegram", "serve"]
# ----------------------------------------------------------