# Telegram Module

Приложение для работы с Telegram

## 🎯 Возможности

- ✅ Обработка webhooks из Telegram ([Telegram Bot API](https://core.telegram.org/bots/api))
- ✅ Постоянное хранилище данных (PostgreSQL)
- ✅ Clean Architecture и CQRS паттерны

## 🚀 Локальная раработка

TODO: описать, как работать с проектом локально

## 🏗️ Архитектура

Проект организован согласно **Clean Architecture** с разделением на слои:

```plain

telegram-module/
├── cmd/api/                       # Entry points
├── internal/
│   ├── domain/                    # DDD: Entities, Value Objects, Errors
│   ├── application/               # CQRS: Commands, Queries, DTOs
│   ├── infrastructure/            # Telegram API, PostgreSQL, Logger
│   └── presentation/              # Adapters: HTTP, CLI
└── database/postgres/migrations   # Postgres миграции БД

```

### Слои архитектуры

#### 1. **Domain Layer** (`internal/domain/`)

Содержит бизнес-логику приложения, полностью независимую от фреймворков:
- **Entities**: BotWebhookUpdate
- **Domain Errors**: Специфичные для предметной области ошибки

#### 2. **Application Layer** (`internal/application/`)
CQRS реализация:
- **Commands**: SavaWebhookBotUpdate — операция сохранения вэбхуки
- **Quries**: на текущий момент не реализовано

#### 3. **Infrastructure Layer** (`internal/infrastructure/`)
Реализация технических деталей:
- **Telegram Bot Api Service**: Интеграция с Telegram Bot API
- **Storage PostgreSQL Repository**: Слой доступа к данным
- **Config**: Управление конфигурацией приложения

#### 4. **Presentation Layer** (`internal/presentation/`)
Адаптеры входных точек:
- **HTTP Handler**: REST API для web клиентов

## 🔧 Конфигурация

### Переменные окружения

Создайте `.env` файл на основе `.env.example`:

### Получение Telegram credentials

Для получения токена для работы с ботом, используйте @BotFather

## 📊 База данных

Приложение работает с PostgreSQL с использованием пакета `sqlx`, смотри [основную документацию](https://github.com/jmoiron/sqlx).

Пример работы для получения данных

```go

type User struct {
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    Uuid string
}

// Query the database, storing results in a []User (wrapped in []interface{})
user := []User{}
db.Select(&user, "SELECT first_name, last_name, uuid FROM users ORDER BY first_name ASC")
fmt.Printf("%#v", user)

```

### Миграции

Миграции находятся в `db/postgresql/migrations`.

При работе с миграциями используем [пакет](https://github.com/golang-migrate/migrate).

Для локальной работы с миграциями - использовать [CLI](https://github.com/golang-migrate/migrate/tree/v4.19.1/cmd/migrate).

Пример создания миграции

```bash
migrate create -ext sql -dir ./db/postgres/migrations create_table_webhooks
```

Пример запуска миграции

```bash
migrate -source file://./database/postgres/migrations -database "postgresql://root:qwerty123@localhost:5432/telegram?sslmode=disable" up 
```

## 🧪 Тестирование

### Запуск всех тестов

TODO: необходимо описать тесты

### Покрытие тестами

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔒 Стандарты кода

### Go Conventions

TODO: необходимо вначале понять общепринятые и потом сюда описать

### Комментарии

```go
// Package entities содержит основные бизнес-сущности
package entities

// Message представляет сообщение из Telegram канала
type Message struct {
    // поля
}

// NewMessage создает новое сообщение с валидацией
func NewMessage(...) (*Message, error) {
    // реализация
}
```

## 🚨 Обработка ошибок

Приложение использует следующую иерархию ошибок:

```plain
Domain Errors (специфичные для бизнес-логики)
  ├── ErrInvalidChannelID
  ├── ErrInvalidMessage
  └── ErrNoMessages

Infrastructure Errors (ошибки от внешних сервисов)
  ├── ErrDatabaseConnection
  └── ErrTelegramAPI

Application Errors (ошибки применения)
  └── Оборачивают domain и infrastructure ошибки
```
