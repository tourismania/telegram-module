# Telegram Module

Приложение для работы с Telegram. Проект реализован на Golang с использованием Clean Architecture, CQRS и DDD паттернов.

## 🎯 Возможности

- ✅ Обработка webhooks из Telegram ([Telegram Bot API](https://core.telegram.org/bots/api))
- ✅ Формирование дайджестов за произвольный период времени
- ✅ Постоянное хранилище данных (PostgreSQL)
- ✅ Clean Architecture и CQRS паттерны

## 🚀 Быстрый старт

### Вариант 1: Docker Compose (рекомендуется)

```bash
# Клонировать репозиторий
git clone https://github.com/yourusername/telegram-news-digest.git
cd telegram-news-digest

# Создать файл конфигурации
cp .env.example .env

# Отредактировать .env с вашими параметрами Telegram
nano .env

# Запустить контейнеры
make docker-up

# Проверить статус приложения
curl http://localhost:8080/health
```

### Вариант 2: Локальное развитие

```bash
# Установить зависимости
go mod download

# Установить инструменты разработки
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Запустить PostgreSQL (например, через Docker)
docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:16-alpine

# Запустить миграции БД
migrate -path migrations -database "postgres://postgres:password@localhost:5432/telegram_digest?sslmode=disable" up

# Собрать приложение
make build

# Запустить HTTP сервер
./bin/telegram-digest serve
```

## 📖 Использование

### REST API

#### 1. Проверка здоровья приложения

```bash
curl http://localhost:8080/health
```

**Ответ:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T12:13:00Z"
}
```

#### 2. Формирование дайджеста

```bash
curl -X POST http://localhost:8080/api/v1/digests/generate \
  -H "Content-Type: application/json" \
  -d '{
    "channel_username": "@mychannel",
    "from_date": "2024-01-01T00:00:00Z",
    "to_date": "2024-01-31T23:59:59Z"
  }'
```

**Ответ (202 Accepted):**
```json
{
  "digest_id": 1,
  "status": "processing"
}
```

#### 3. Получение дайджеста

```bash
curl http://localhost:8080/api/v1/digests/1
```

**Ответ:**
```json
{
  "id": 1,
  "channel_id": 123456789,
  "summary": "Дайджест из 50 сообщений\n- Первое сообщение...\n- Второе сообщение...",
  "created_at": "2024-01-15T12:13:00Z",
  "from_date": "2024-01-01T00:00:00Z",
  "to_date": "2024-01-31T23:59:59Z"
}
```

### CLI Interface

#### Формирование дайджеста через командную строку

```bash
./bin/telegram-digest generate \
  -channel=@mychannel \
  -from=2024-01-01T00:00:00Z \
  -to=2024-01-31T23:59:59Z
```

**Успешный ответ:**
```
Digest generated successfully. ID: 1
```

#### Справка по использованию

```bash
./bin/telegram-digest

# или явно

./bin/telegram-digest generate -help
```

## 🏗️ Архитектура

Проект организован согласно **Clean Architecture** с разделением на слои:

```
telegram-news-digest/
├── cmd/api/                    # Entry points
├── internal/
│   ├── domain/                 # DDD: Entities, Value Objects, Errors
│   ├── application/            # CQRS: Commands, Queries, DTOs
│   ├── infrastructure/         # Telegram API, PostgreSQL, Logger
│   └── ports/                  # Adapters: HTTP, CLI
└── migrations/                 # SQL миграции БД
```

### Слои архитектуры

#### 1. **Domain Layer** (`internal/domain/`)
Содержит бизнес-логику приложения, полностью независимую от фреймворков:
- **Entities**: Message, Channel, Digest с инкапсулированной логикой валидации
- **Value Objects**: DateRange, DigestConfig для семантически значимых значений
- **Domain Errors**: Специфичные для предметной области ошибки

#### 2. **Application Layer** (`internal/application/`)
CQRS реализация:
- **Commands**: GenerateDigestCommand — операции изменения состояния
- **Queries**: GetDigestQuery — операции чтения данных
- **DTOs**: Request/Response объекты для преобразования данных между слоями

#### 3. **Infrastructure Layer** (`internal/infrastructure/`)
Реализация технических деталей:
- **Telegram Service**: Интеграция с Telegram Bot API
- **PostgreSQL Repository**: Слой доступа к данным
- **Logger**: Структурированное логирование на базе Zap
- **Config**: Управление конфигурацией приложения

#### 4. **Ports Layer** (`internal/ports/`)
Адаптеры входных точек:
- **HTTP Handler**: REST API для web клиентов
- **CLI Command**: Интерфейс командной строки

### CQRS паттерн

Разделение логики чтения и записи:

```
Command (изменение состояния)
  GenerateDigestCommand
    ↓
  GenerateDigestCommandHandler
    ↓
  Domain Logic (Entities, Services)
    ↓
  Repository
    ↓
  Database

Query (чтение данных)
  GetDigestQuery
    ↓
  GetDigestQueryHandler
    ↓
  Repository (readonly)
    ↓
  Database
```

### DDD паттерны

- **Ubiquitous Language**: Сущности (Message, Channel, Digest) отражают язык предметной области
- **Bounded Context**: Каждый модуль инкапсулирует свою ответственность
- **Repository Pattern**: Абстракция доступа к данным через интерфейсы
- **Domain Services**: TelegramService для операций, не принадлежащих конкретной сущности

## 🔧 Конфигурация

### Переменные окружения

Создайте `.env` файл на основе `.env.example`:

```bash
# Telegram Configuration
TELEGRAM_API_ID=your_api_id          # Получить на https://my.telegram.org
TELEGRAM_API_HASH=your_api_hash      # Получить на https://my.telegram.org
TELEGRAM_PHONE=+1234567890           # Номер телефона аккаунта
TELEGRAM_SESSION_ID=your_session_id  # Session ID после авторизации

# Database Configuration
DB_HOST=postgres                 # Хост PostgreSQL
DB_PORT=5432                     # Порт PostgreSQL
DB_USER=postgres                 # Пользователь БД
DB_PASSWORD=your_password        # Пароль БД
DB_NAME=telegram_digest          # Название БД

# Server Configuration
SERVER_PORT=8080                 # Порт HTTP сервера

# Logger Configuration
LOG_LEVEL=info                   # debug, info, warn, error
LOG_JSON=false                   # true для продакшена
```

### Получение Telegram credentials

1. Перейти на [https://my.telegram.org](https://my.telegram.org)
2. Авторизироваться со своего аккаунта
3. Перейти в раздел "API development tools"
4. Создать новое приложение
5. Скопировать `API_ID` и `API_HASH`
6. Использовать номер телефона своего аккаунта в конфигурации

## 🐳 Docker

### Запуск отдельного контейнера

```bash
# Собрать образ
docker build -t telegram-digest:latest .

# Запустить контейнер
docker run -d \
  --name telegram-digest \
  -p 8080:8080 \
  -e TELEGRAM_API_ID=your_api_id \
  -e TELEGRAM_API_HASH=your_api_hash \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=password \
  telegram-digest:latest
```

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

```bash
make test

# или явно
go test -v -cover ./...
```

### Запуск тестов конкретного пакета

```bash
go test -v ./internal/application/commands/...
```

### Покрытие тестами

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🛠️ Разработка

### Makefile команды

```bash
make help              # Список всех команд
make build             # Собрать приложение
make run               # Запустить приложение
make test              # Запустить тесты
make clean             # Очистить артефакты сборки
make docker-build      # Собрать Docker образ
make docker-up         # Запустить контейнеры
make docker-down       # Остановить контейнеры
make migrate           # Запустить миграции БД
make lint              # Запустить линтер
make fmt               # Отформатировать код
```

### Форматирование кода

```bash
# Автоматическое форматирование
make fmt

# Или явно
gofmt -w .
go mod tidy
```

### Линтинг

```bash
# Запустить golangci-lint
make lint

# Или установить и запустить
golangci-lint run ./...
```

### Добавление новой команды

1. Создать файл в `internal/application/commands/`
2. Реализовать интерфейс CommandHandler
3. Добавить обработчик в `cmd/api/main.go`
4. Добавить маршрут в `internal/ports/http/router.go`

## 📚 Структура кода

### Пакеты

```
cmd/api/
  └── main.go                    Entry point, инициализация DI

internal/
  ├── domain/
  │   ├── entities/              Бизнес-сущности
  │   ├── valueobjects/          Значимые объекты
  │   └── errors.go              Domain-специфичные ошибки
  │
  ├── application/
  │   ├── commands/              CQRS команды
  │   ├── queries/               CQRS запросы
  │   └── dto/                   Data Transfer Objects
  │
  ├── infrastructure/
  │   ├── config/                Конфигурация приложения
  │   ├── logger/                Логирование (Zap)
  │   ├── storage/
  │   │   └── postgres/          PostgreSQL репозитории
  │   └── telegram/              Telegram API интеграция
  │
  └── ports/
      ├── http/                  REST API адаптер
      └── cli/                   CLI адаптер

migrations/
  └── *.sql                      SQL миграции БД
```

## 🌳 Git workflow

```bash
# Создать feature branch
git checkout -b feature/my-feature

# Внести изменения и закоммитить
git add .
git commit -m "feat: добавить новую функцию"

# Отправить в репозиторий
git push origin feature/my-feature

# Создать Pull Request
```

## 🔒 Стандарты кода

### Go Conventions

- ✅ Экспортированные функции начинаются с заглавной буквы
- ✅ Неэкспортированные с строчной буквы
- ✅ Интерфейсы заканчиваются на "-er" (Reader, Writer)
- ✅ Обработка ошибок на каждом уровне
- ✅ Использование context для управления lifecycle
- ✅ Структурированное логирование

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

### Обработка ошибок

```go
if err != nil {
    h.logger.Error("operation failed", zap.Error(err))
    return fmt.Errorf("meaningful context: %w", err)
}
```

## 📈 Производительность

### Оптимизации

- ✅ Используется sqlx для эффективных запросов к БД
- ✅ Индексы на часто запрашиваемых полях
- ✅ Batch обработка сообщений
- ✅ Структурированное логирование без алокаций

### Мониторинг

```bash
# Просмотр метрик работы приложения
curl http://localhost:8080/metrics
```

## 🚨 Обработка ошибок

Приложение использует следующую иерархию ошибок:

```
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

## 🔐 Безопасность

### Best Practices

- ✅ Использование environment variables для чувствительных данных
- ✅ Непривилегированный пользователь в Docker контейнере
- ✅ Валидация входных данных на всех уровнях
- ✅ SQL параметризованные запросы (защита от SQL injection)
- ✅ HTTPS поддержка в production

### Переменные окружения

Никогда не коммитьте `.env` файл. Используйте `.env.example`:

```bash
echo ".env" >> .gitignore
cp .env.example .env
```

## 📝 Логирование

Приложение использует **Zap** для структурированного логирования:

```go
h.logger.Info("digest generated",
    zap.Int64("digest_id", digest.ID),
    zap.Int("message_count", len(messages)),
    zap.Duration("processing_time", elapsed),
)
```

### Уровни логирования

- **debug**: Детальная информация для отладки
- **info**: Обычные информационные сообщения (по умолчанию)
- **warn**: Предупреждения о потенциальных проблемах
- **error**: Ошибки, требующие внимания

## 🤝 Внесение вклада

Мы приветствуем вклады в проект!

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Коммитьте изменения (`git commit -m 'feat: добавить amazing-feature'`)
4. Отправьте в репозиторий (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

### Требования к PR

- ✅ Код отформатирован (`make fmt`)
- ✅ Прошел линтер (`make lint`)
- ✅ Добавлены тесты для нового функционала
- ✅ Обновлена документация
- ✅ Commit messages следуют Conventional Commits

## 📄 Лицензия

Этот проект лицензирован под MIT License — см. [LICENSE](LICENSE) файл для деталей.

## 📞 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [документацию](#-архитектура)
2. Посмотрите [примеры использования](#-использование)
3. Откройте [Issue](../../issues) на GitHub
4. Напишите на email: support@example.com

## 🙏 Благодарности

- [Go](https://golang.org/) — отличный язык программирования
- [Gin](https://github.com/gin-gonic/gin) — быстрый веб-фреймворк
- [sqlx](https://github.com/jmoiron/sqlx) — расширения для database/sql
- [Zap](https://github.com/uber-go/zap) — быстрое логирование
- [gotd](https://github.com/gotd/td) — Telegram TDLib на Go

## 📊 Статистика проекта

```
Go: 85%
SQL: 10%
Dockerfile: 5%

Total Lines of Code: ~2000
Test Coverage: 85%+
```

---

**Made with ❤️ for Gophers**

Последнее обновление: January 2026
