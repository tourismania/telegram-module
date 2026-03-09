# Проект Telegram News Digest - Архитектура и Реализация

## 1. Структура проекта

```
telegram-news-digest/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point приложения
├── internal/
│   ├── domain/                        # DDD Domain Layer
│   │   ├── entities/
│   │   │   ├── message.go
│   │   │   ├── channel.go
│   │   │   └── digest.go
│   │   ├── valueobjects/
│   │   │   ├── date_range.go
│   │   │   └── digest_config.go
│   │   └── errors.go
│   ├── application/                   # Application Layer (CQRS)
│   │   ├── commands/
│   │   │   ├── generate_digest.go
│   │   │   └── handler.go
│   │   ├── queries/
│   │   │   ├── get_digest.go
│   │   │   └── handler.go
│   │   └── dto/
│   │       ├── request.go
│   │       └── response.go
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── telegram/
│   │   │   ├── client.go
│   │   │   └── service.go
│   │   ├── storage/
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go
│   │   │   │   ├── migrations/
│   │   │   │   └── repository.go
│   │   │   └── interfaces.go
│   │   ├── logger/
│   │   │   └── zap_logger.go
│   │   └── config/
│   │       └── config.go
│   └── ports/                        # Adapter Layer
│       ├── http/
│       │   ├── handler.go
│       │   └── router.go
│       └── cli/
│           └── command.go
├── migrations/
│   └── 001_initial_schema.sql
├── config/
│   └── config.yaml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md
```

## 2. Реализация компонентов

### 2.1 Domain Layer

**internal/domain/entities/message.go**
```go
package entities

import (
	"time"
)

// Message представляет сообщение из Telegram канала
type Message struct {
	ID        int64
	ChannelID int64
	Text      string
	MediaURL  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMessage создает новое сообщение с валидацией
func NewMessage(channelID int64, text string, createdAt time.Time) (*Message, error) {
	if channelID <= 0 {
		return nil, ErrInvalidChannelID
	}
	if text == "" {
		return nil, ErrInvalidMessage
	}

	return &Message{
		ChannelID: channelID,
		Text:      text,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}, nil
}
```

**internal/domain/entities/channel.go**
```go
package entities

// Channel представляет Telegram канал
type Channel struct {
	ID          int64
	Username    string
	Title       string
	Description string
	AccessHash  int64
}

// NewChannel создает новый канал
func NewChannel(id int64, username string, title string) (*Channel, error) {
	if id <= 0 {
		return nil, ErrInvalidChannelID
	}
	if username == "" {
		return nil, ErrInvalidChannelUsername
	}

	return &Channel{
		ID:       id,
		Username: username,
		Title:    title,
	}, nil
}
```

**internal/domain/entities/digest.go**
```go
package entities

import (
	"time"
)

// Digest представляет дайджест новостей
type Digest struct {
	ID        int64
	ChannelID int64
	Messages  []Message
	Summary   string
	CreatedAt time.Time
	FromDate  time.Time
	ToDate    time.Time
}

// NewDigest создает новый дайджест
func NewDigest(channelID int64, messages []Message, fromDate, toDate time.Time) (*Digest, error) {
	if channelID <= 0 {
		return nil, ErrInvalidChannelID
	}
	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}
	if len(messages) == 0 {
		return nil, ErrNoMessages
	}

	return &Digest{
		ChannelID: channelID,
		Messages:  messages,
		CreatedAt: time.Now(),
		FromDate:  fromDate,
		ToDate:    toDate,
	}, nil
}

// AddSummary добавляет сводку к дайджесту
func (d *Digest) AddSummary(summary string) error {
	if summary == "" {
		return ErrInvalidSummary
	}
	d.Summary = summary
	return nil
}
```

**internal/domain/valueobjects/date_range.go**
```go
package valueobjects

import (
	"time"
)

// DateRange представляет диапазон дат
type DateRange struct {
	Start time.Time
	End   time.Time
}

// NewDateRange создает новый диапазон дат с валидацией
func NewDateRange(start, end time.Time) (*DateRange, error) {
	if start.After(end) {
		return nil, ErrInvalidDateRange
	}
	return &DateRange{
		Start: start,
		End:   end,
	}, nil
}

// IsInRange проверяет, находится ли дата в диапазоне
func (dr *DateRange) IsInRange(date time.Time) bool {
	return !date.Before(dr.Start) && !date.After(dr.End)
}
```

**internal/domain/errors.go**
```go
package domain

import "errors"

var (
	ErrInvalidChannelID       = errors.New("invalid channel id")
	ErrInvalidChannelUsername = errors.New("invalid channel username")
	ErrInvalidMessage         = errors.New("invalid message")
	ErrInvalidDateRange       = errors.New("invalid date range: start must be before end")
	ErrNoMessages             = errors.New("no messages found")
	ErrInvalidSummary         = errors.New("invalid summary")
	ErrChannelNotFound        = errors.New("channel not found")
	ErrDigestNotFound         = errors.New("digest not found")
)
```

### 2.2 Infrastructure Layer

**internal/infrastructure/config/config.go**
```go
package config

import (
	"os"
	"strconv"
	"time"
)

// Config содержит конфигурацию приложения
type Config struct {
	Telegram TelegramConfig
	Database DatabaseConfig
	Server   ServerConfig
	Logger   LoggerConfig
}

type TelegramConfig struct {
	APIHash   string
	APIKey    string
	Phone     string
	SessionID string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LoggerConfig struct {
	Level string
	JSON  bool
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() Config {
	return Config{
		Telegram: TelegramConfig{
			APIHash:   os.Getenv("TELEGRAM_API_HASH"),
			APIKey:    os.Getenv("TELEGRAM_API_ID"),
			Phone:     os.Getenv("TELEGRAM_PHONE"),
			SessionID: os.Getenv("TELEGRAM_SESSION_ID"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   getEnv("DB_NAME", "telegram_digest"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 15)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 15)) * time.Second,
			IdleTimeout:  time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			JSON:  getEnvBool("LOG_JSON", false),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvBool(name string, defaultVal bool) bool {
	val := getEnv(name, "")
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1"
}
```

**internal/infrastructure/logger/zap_logger.go**
```go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger интерфейс логирования
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type zapLogger struct {
	log *zap.Logger
}

// NewLogger создает новый logger
func NewLogger(level string, json bool) (Logger, error) {
	config := zap.NewDevelopmentConfig()
	if json {
		config = zap.NewProductionConfig()
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{log: log}, nil
}

func (zl *zapLogger) Debug(msg string, fields ...zap.Field) {
	zl.log.Debug(msg, fields...)
}

func (zl *zapLogger) Info(msg string, fields ...zap.Field) {
	zl.log.Info(msg, fields...)
}

func (zl *zapLogger) Warn(msg string, fields ...zap.Field) {
	zl.log.Warn(msg, fields...)
}

func (zl *zapLogger) Error(msg string, fields ...zap.Field) {
	zl.log.Error(msg, fields...)
}

func (zl *zapLogger) Fatal(msg string, fields ...zap.Field) {
	zl.log.Fatal(msg, fields...)
}
```

**internal/infrastructure/storage/interfaces.go**
```go
package storage

import (
	"context"
	"time"

	"telegram-news-digest/internal/domain/entities"
)

// MessageRepository интерфейс для работы с сообщениями
type MessageRepository interface {
	Save(ctx context.Context, message *entities.Message) error
	FindByChannelAndDateRange(ctx context.Context, channelID int64, from, to time.Time) ([]entities.Message, error)
	FindByID(ctx context.Context, id int64) (*entities.Message, error)
}

// ChannelRepository интерфейс для работы с каналами
type ChannelRepository interface {
	Save(ctx context.Context, channel *entities.Channel) error
	FindByID(ctx context.Context, id int64) (*entities.Channel, error)
	FindByUsername(ctx context.Context, username string) (*entities.Channel, error)
}

// DigestRepository интерфейс для работы с дайджестами
type DigestRepository interface {
	Save(ctx context.Context, digest *entities.Digest) error
	FindByID(ctx context.Context, id int64) (*entities.Digest, error)
	FindByChannelAndDateRange(ctx context.Context, channelID int64, from, to time.Time) ([]entities.Digest, error)
}
```

**internal/infrastructure/storage/postgres/connection.go**
```go
package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"telegram-news-digest/internal/infrastructure/config"
)

// NewConnection создает новое соединение с PostgreSQL
func NewConnection(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
```

**internal/infrastructure/storage/postgres/repository.go**
```go
package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"telegram-news-digest/internal/domain/entities"
)

// MessageRepository реализация репозитория сообщений
type MessageRepository struct {
	db *sqlx.DB
}

// NewMessageRepository создает новый репозиторий сообщений
func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Save сохраняет сообщение в БД
func (r *MessageRepository) Save(ctx context.Context, message *entities.Message) error {
	query := `
		INSERT INTO messages (channel_id, text, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		message.ChannelID,
		message.Text,
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&message.ID)

	return err
}

// FindByChannelAndDateRange находит сообщения по каналу и диапазону дат
func (r *MessageRepository) FindByChannelAndDateRange(
	ctx context.Context,
	channelID int64,
	from, to time.Time,
) ([]entities.Message, error) {
	var messages []entities.Message

	query := `
		SELECT id, channel_id, text, created_at, updated_at
		FROM messages
		WHERE channel_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &messages, query, channelID, from, to)
	return messages, err
}

// FindByID находит сообщение по ID
func (r *MessageRepository) FindByID(ctx context.Context, id int64) (*entities.Message, error) {
	var message entities.Message

	query := `
		SELECT id, channel_id, text, created_at, updated_at
		FROM messages WHERE id = $1
	`

	err := r.db.GetContext(ctx, &message, query, id)
	return &message, err
}

// ChannelRepository реализация репозитория каналов
type ChannelRepository struct {
	db *sqlx.DB
}

// NewChannelRepository создает новый репозиторий каналов
func NewChannelRepository(db *sqlx.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

// Save сохраняет канал в БД
func (r *ChannelRepository) Save(ctx context.Context, channel *entities.Channel) error {
	query := `
		INSERT INTO channels (id, username, title, description, access_hash)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username, title = EXCLUDED.title
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		channel.ID,
		channel.Username,
		channel.Title,
		channel.Description,
		channel.AccessHash,
	)

	return err
}

// FindByID находит канал по ID
func (r *ChannelRepository) FindByID(ctx context.Context, id int64) (*entities.Channel, error) {
	var channel entities.Channel

	query := `
		SELECT id, username, title, description, access_hash
		FROM channels WHERE id = $1
	`

	err := r.db.GetContext(ctx, &channel, query, id)
	return &channel, err
}

// FindByUsername находит канал по username
func (r *ChannelRepository) FindByUsername(ctx context.Context, username string) (*entities.Channel, error) {
	var channel entities.Channel

	query := `
		SELECT id, username, title, description, access_hash
		FROM channels WHERE username = $1
	`

	err := r.db.GetContext(ctx, &channel, query, username)
	return &channel, err
}

// DigestRepository реализация репозитория дайджестов
type DigestRepository struct {
	db *sqlx.DB
}

// NewDigestRepository создает новый репозиторий дайджестов
func NewDigestRepository(db *sqlx.DB) *DigestRepository {
	return &DigestRepository{db: db}
}

// Save сохраняет дайджест в БД
func (r *DigestRepository) Save(ctx context.Context, digest *entities.Digest) error {
	query := `
		INSERT INTO digests (channel_id, summary, created_at, from_date, to_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		digest.ChannelID,
		digest.Summary,
		digest.CreatedAt,
		digest.FromDate,
		digest.ToDate,
	).Scan(&digest.ID)

	return err
}

// FindByID находит дайджест по ID
func (r *DigestRepository) FindByID(ctx context.Context, id int64) (*entities.Digest, error) {
	var digest entities.Digest

	query := `
		SELECT id, channel_id, summary, created_at, from_date, to_date
		FROM digests WHERE id = $1
	`

	err := r.db.GetContext(ctx, &digest, query, id)
	return &digest, err
}

// FindByChannelAndDateRange находит дайджесты по каналу и диапазону дат
func (r *DigestRepository) FindByChannelAndDateRange(
	ctx context.Context,
	channelID int64,
	from, to time.Time,
) ([]entities.Digest, error) {
	var digests []entities.Digest

	query := `
		SELECT id, channel_id, summary, created_at, from_date, to_date
		FROM digests
		WHERE channel_id = $1 AND from_date >= $2 AND to_date <= $3
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &digests, query, channelID, from, to)
	return digests, err
}
```

**internal/infrastructure/storage/postgres/migrations/001_initial_schema.sql**
```sql
CREATE TABLE channels (
    id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255),
    description TEXT,
    access_hash BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    media_url TEXT[],
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_channel_created ON messages(channel_id, created_at DESC);

CREATE TABLE digests (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    from_date TIMESTAMP NOT NULL,
    to_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_digests_channel_dates ON digests(channel_id, from_date, to_date);
```

**internal/infrastructure/telegram/service.go**
```go
package telegram

import (
	"context"
	"time"

	"github.com/gotd/td/tdjson"
	"go.uber.org/zap"

	"telegram-news-digest/internal/domain/entities"
	"telegram-news-digest/internal/infrastructure/logger"
)

// TelegramService сервис для работы с Telegram API
type TelegramService struct {
	client tdjson.ServerBridge
	logger logger.Logger
}

// NewTelegramService создает новый сервис Telegram
func NewTelegramService(logger logger.Logger) *TelegramService {
	return &TelegramService{
		logger: logger,
	}
}

// FetchMessages получает сообщения из канала за период времени
func (ts *TelegramService) FetchMessages(
	ctx context.Context,
	channelID int64,
	from, to time.Time,
) ([]entities.Message, error) {
	ts.logger.Info("fetching messages from channel", 
		zap.Int64("channel_id", channelID),
		zap.Time("from", from),
		zap.Time("to", to),
	)

	// Реализация зависит от используемой Telegram библиотеки (gotd, telego, т.д.)
	// Здесь представлен примерный паттерн

	var messages []entities.Message

	// Логика получения сообщений из API
	ts.logger.Info("messages fetched successfully", zap.Int("count", len(messages)))

	return messages, nil
}

// GetChannelInfo получает информацию о канале
func (ts *TelegramService) GetChannelInfo(ctx context.Context, username string) (*entities.Channel, error) {
	ts.logger.Info("fetching channel info", zap.String("username", username))

	// Реализация получения информации о канале

	return &entities.Channel{
		Username: username,
	}, nil
}
```

### 2.3 Application Layer (CQRS)

**internal/application/dto/request.go**
```go
package dto

import "time"

// GenerateDigestRequest DTO для запроса на формирование дайджеста
type GenerateDigestRequest struct {
	ChannelUsername string    `json:"channel_username" binding:"required"`
	FromDate        time.Time `json:"from_date" binding:"required"`
	ToDate          time.Time `json:"to_date" binding:"required"`
}

// GetDigestRequest DTO для получения дайджеста
type GetDigestRequest struct {
	DigestID int64 `uri:"id" binding:"required"`
}
```

**internal/application/dto/response.go**
```go
package dto

import "time"

// DigestResponse DTO для ответа с дайджестом
type DigestResponse struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
}

// MessageResponse DTO для ответа с сообщением
type MessageResponse struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse стандартный ответ об ошибке
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
```

**internal/application/commands/generate_digest.go**
```go
package commands

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"telegram-news-digest/internal/domain"
	"telegram-news-digest/internal/domain/entities"
	"telegram-news-digest/internal/infrastructure/logger"
	"telegram-news-digest/internal/infrastructure/storage"
	"telegram-news-digest/internal/infrastructure/telegram"
)

// GenerateDigestCommand команда для формирования дайджеста
type GenerateDigestCommand struct {
	ChannelUsername string
	FromDate        time.Time
	ToDate          time.Time
}

// GenerateDigestCommandHandler обработчик команды GenerateDigest
type GenerateDigestCommandHandler struct {
	telegramSvc      *telegram.TelegramService
	messageRepo      storage.MessageRepository
	channelRepo      storage.ChannelRepository
	digestRepo       storage.DigestRepository
	logger           logger.Logger
}

// NewGenerateDigestCommandHandler создает новый обработчик команды
func NewGenerateDigestCommandHandler(
	telegramSvc *telegram.TelegramService,
	messageRepo storage.MessageRepository,
	channelRepo storage.ChannelRepository,
	digestRepo storage.DigestRepository,
	logger logger.Logger,
) *GenerateDigestCommandHandler {
	return &GenerateDigestCommandHandler{
		telegramSvc: telegramSvc,
		messageRepo: messageRepo,
		channelRepo: channelRepo,
		digestRepo:  digestRepo,
		logger:      logger,
	}
}

// Handle обрабатывает команду GenerateDigest
func (h *GenerateDigestCommandHandler) Handle(ctx context.Context, cmd GenerateDigestCommand) (int64, error) {
	h.logger.Info("handling GenerateDigestCommand",
		zap.String("channel_username", cmd.ChannelUsername),
		zap.Time("from_date", cmd.FromDate),
		zap.Time("to_date", cmd.ToDate),
	)

	// 1. Получить информацию о канале
	channel, err := h.telegramSvc.GetChannelInfo(ctx, cmd.ChannelUsername)
	if err != nil {
		h.logger.Error("failed to get channel info", zap.Error(err))
		return 0, fmt.Errorf("failed to get channel info: %w", err)
	}

	// 2. Сохранить канал в репозиторий
	if err := h.channelRepo.Save(ctx, channel); err != nil {
		h.logger.Error("failed to save channel", zap.Error(err))
		return 0, fmt.Errorf("failed to save channel: %w", err)
	}

	// 3. Получить сообщения из Telegram
	messages, err := h.telegramSvc.FetchMessages(ctx, channel.ID, cmd.FromDate, cmd.ToDate)
	if err != nil {
		h.logger.Error("failed to fetch messages", zap.Error(err))
		return 0, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// 4. Создать дайджест
	digest, err := entities.NewDigest(channel.ID, messages, cmd.FromDate, cmd.ToDate)
	if err != nil {
		if err == domain.ErrNoMessages {
			h.logger.Warn("no messages found for digest")
			return 0, err
		}
		h.logger.Error("failed to create digest", zap.Error(err))
		return 0, fmt.Errorf("failed to create digest: %w", err)
	}

	// 5. Сгенерировать сводку (можно использовать AI API)
	summary := h.generateSummary(messages)
	if err := digest.AddSummary(summary); err != nil {
		h.logger.Error("failed to add summary", zap.Error(err))
		return 0, fmt.Errorf("failed to add summary: %w", err)
	}

	// 6. Сохранить все сообщения
	for _, msg := range messages {
		msg.ChannelID = channel.ID
		if err := h.messageRepo.Save(ctx, &msg); err != nil {
			h.logger.Error("failed to save message", zap.Error(err))
			// Продолжаем, несмотря на ошибку
		}
	}

	// 7. Сохранить дайджест
	if err := h.digestRepo.Save(ctx, digest); err != nil {
		h.logger.Error("failed to save digest", zap.Error(err))
		return 0, fmt.Errorf("failed to save digest: %w", err)
	}

	h.logger.Info("digest generated successfully", zap.Int64("digest_id", digest.ID))
	return digest.ID, nil
}

// generateSummary генерирует сводку из сообщений
func (h *GenerateDigestCommandHandler) generateSummary(messages []entities.Message) string {
	// Простая реализация: объединение первых N сообщений
	// В реальном приложении можно использовать AI сервис
	summary := fmt.Sprintf("Дайджест из %d сообщений\n", len(messages))
	for i, msg := range messages {
		if i >= 5 { // Ограничиваем 5 сообщениями
			break
		}
		if len(msg.Text) > 100 {
			summary += fmt.Sprintf("- %s...\n", msg.Text[:100])
		} else {
			summary += fmt.Sprintf("- %s\n", msg.Text)
		}
	}
	return summary
}
```

**internal/application/queries/get_digest.go**
```go
package queries

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"telegram-news-digest/internal/application/dto"
	"telegram-news-digest/internal/infrastructure/logger"
	"telegram-news-digest/internal/infrastructure/storage"
)

// GetDigestQuery запрос для получения дайджеста
type GetDigestQuery struct {
	DigestID int64
}

// GetDigestQueryHandler обработчик запроса GetDigest
type GetDigestQueryHandler struct {
	digestRepo storage.DigestRepository
	logger     logger.Logger
}

// NewGetDigestQueryHandler создает новый обработчик запроса
func NewGetDigestQueryHandler(
	digestRepo storage.DigestRepository,
	logger logger.Logger,
) *GetDigestQueryHandler {
	return &GetDigestQueryHandler{
		digestRepo: digestRepo,
		logger:     logger,
	}
}

// Handle обрабатывает запрос GetDigest
func (h *GetDigestQueryHandler) Handle(ctx context.Context, query GetDigestQuery) (*dto.DigestResponse, error) {
	h.logger.Info("handling GetDigestQuery", zap.Int64("digest_id", query.DigestID))

	digest, err := h.digestRepo.FindByID(ctx, query.DigestID)
	if err != nil {
		h.logger.Error("failed to find digest", zap.Error(err))
		return nil, fmt.Errorf("failed to find digest: %w", err)
	}

	return &dto.DigestResponse{
		ID:        digest.ID,
		ChannelID: digest.ChannelID,
		Summary:   digest.Summary,
		CreatedAt: digest.CreatedAt,
		FromDate:  digest.FromDate,
		ToDate:    digest.ToDate,
	}, nil
}
```

### 2.4 Adapter Layer (HTTP)

**internal/ports/http/handler.go**
```go
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"telegram-news-digest/internal/application/commands"
	"telegram-news-digest/internal/application/dto"
	"telegram-news-digest/internal/application/queries"
	"telegram-news-digest/internal/infrastructure/logger"
)

// DigestHandler обработчик HTTP запросов для дайджестов
type DigestHandler struct {
	generateDigestCmdHandler *commands.GenerateDigestCommandHandler
	getDigestQueryHandler    *queries.GetDigestQueryHandler
	logger                   logger.Logger
}

// NewDigestHandler создает новый HTTP обработчик
func NewDigestHandler(
	generateDigestCmdHandler *commands.GenerateDigestCommandHandler,
	getDigestQueryHandler *queries.GetDigestQueryHandler,
	logger logger.Logger,
) *DigestHandler {
	return &DigestHandler{
		generateDigestCmdHandler: generateDigestCmdHandler,
		getDigestQueryHandler:    getDigestQueryHandler,
		logger:                   logger,
	}
}

// GenerateDigest генерирует дайджест
// @Summary Generate digest
// @Description Create a news digest for a Telegram channel for a specified period
// @Accept json
// @Produce json
// @Param request body dto.GenerateDigestRequest true "Request body"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/digests/generate [post]
func (h *DigestHandler) GenerateDigest(c *gin.Context) {
	var req dto.GenerateDigestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Валидация дат
	if req.FromDate.After(req.ToDate) {
		h.logger.Warn("invalid date range")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_DATE_RANGE",
			Message: "from_date must be before to_date",
		})
		return
	}

	// Выполнить команду
	cmd := commands.GenerateDigestCommand{
		ChannelUsername: req.ChannelUsername,
		FromDate:        req.FromDate,
		ToDate:          req.ToDate,
	}

	digestID, err := h.generateDigestCmdHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("failed to generate digest", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "DIGEST_GENERATION_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, map[string]interface{}{
		"digest_id": digestID,
		"status":    "processing",
	})
}

// GetDigest получает дайджест по ID
// @Summary Get digest
// @Description Retrieve a news digest by ID
// @Produce json
// @Param id path int true "Digest ID"
// @Success 200 {object} dto.DigestResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/digests/{id} [get]
func (h *DigestHandler) GetDigest(c *gin.Context) {
	var req dto.GetDigestRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Выполнить запрос
	query := queries.GetDigestQuery{
		DigestID: req.DigestID,
	}

	digest, err := h.getDigestQueryHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("failed to get digest", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "DIGEST_NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, digest)
}

// Health проверка здоровья приложения
// @Summary Health check
// @Description Check if the application is running
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *DigestHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
	})
}
```

**internal/ports/http/router.go**
```go
package http

import (
	"github.com/gin-gonic/gin"

	"telegram-news-digest/internal/infrastructure/logger"
)

// Router конфигурирует маршруты приложения
type Router struct {
	engine  *gin.Engine
	handler *DigestHandler
	logger  logger.Logger
}

// NewRouter создает новый роутер
func NewRouter(handler *DigestHandler, logger logger.Logger) *Router {
	return &Router{
		engine:  gin.Default(),
		handler: handler,
		logger:  logger,
	}
}

// Setup настраивает все маршруты
func (r *Router) Setup() {
	// Health check
	r.engine.GET("/health", r.handler.Health)

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		digests := v1.Group("/digests")
		{
			digests.POST("/generate", r.handler.GenerateDigest)
			digests.GET("/:id", r.handler.GetDigest)
		}
	}
}

// Run запускает HTTP сервер
func (r *Router) Run(port string) error {
	r.logger.Info("starting HTTP server on " + port)
	return r.engine.Run(":" + port)
}
```

### 2.5 CLI Interface

**internal/ports/cli/command.go**
```go
package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"telegram-news-digest/internal/application/commands"
	"telegram-news-digest/internal/infrastructure/logger"
)

// CLIApp приложение для работы из командной строки
type CLIApp struct {
	generateDigestCmdHandler *commands.GenerateDigestCommandHandler
	logger                   logger.Logger
}

// NewCLIApp создает новое CLI приложение
func NewCLIApp(
	generateDigestCmdHandler *commands.GenerateDigestCommandHandler,
	logger logger.Logger,
) *CLIApp {
	return &CLIApp{
		generateDigestCmdHandler: generateDigestCmdHandler,
		logger:                   logger,
	}
}

// Run запускает CLI приложение
func (app *CLIApp) Run(args []string) error {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	channelUsername := generateCmd.String("channel", "", "Telegram channel username")
	fromDate := generateCmd.String("from", "", "Start date (RFC3339 format)")
	toDate := generateCmd.String("to", "", "End date (RFC3339 format)")

	if len(args) == 0 {
		app.printUsage()
		return fmt.Errorf("no command provided")
	}

	switch args[0] {
	case "generate":
		if err := generateCmd.Parse(args[1:]); err != nil {
			return err
		}

		if *channelUsername == "" || *fromDate == "" || *toDate == "" {
			fmt.Println("Error: all flags are required")
			generateCmd.PrintDefaults()
			return fmt.Errorf("missing required flags")
		}

		from, err := time.Parse(time.RFC3339, *fromDate)
		if err != nil {
			return fmt.Errorf("invalid from date format: %w", err)
		}

		to, err := time.Parse(time.RFC3339, *toDate)
		if err != nil {
			return fmt.Errorf("invalid to date format: %w", err)
		}

		return app.generateDigest(*channelUsername, from, to)

	default:
		app.printUsage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

// generateDigest генерирует дайджест через CLI
func (app *CLIApp) generateDigest(channel string, from, to time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := commands.GenerateDigestCommand{
		ChannelUsername: channel,
		FromDate:        from,
		ToDate:          to,
	}

	digestID, err := app.generateDigestCmdHandler.Handle(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to generate digest: %w", err)
	}

	log.Printf("Digest generated successfully. ID: %d\n", digestID)
	return nil
}

// printUsage выводит справку по использованию
func (app *CLIApp) printUsage() {
	fmt.Println("Usage: telegram-digest <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  generate  Generate a news digest for a Telegram channel")
	fmt.Println("\nExamples:")
	fmt.Println("  telegram-digest generate -channel=@mychannel -from=2024-01-01T00:00:00Z -to=2024-01-31T23:59:59Z")
}
```

### 2.6 Main Entry Point

**cmd/api/main.go**
```go
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jmoiron/sqlx"

	"telegram-news-digest/internal/application/commands"
	"telegram-news-digest/internal/application/queries"
	"telegram-news-digest/internal/infrastructure/config"
	"telegram-news-digest/internal/infrastructure/logger"
	"telegram-news-digest/internal/infrastructure/storage/postgres"
	"telegram-news-digest/internal/infrastructure/telegram"
	"telegram-news-digest/internal/ports/cli"
	"telegram-news-digest/internal/ports/http"
)

func main() {
	// Загрузить конфигурацию
	cfg := config.LoadConfig()

	// Инициализировать логгер
	log, err := logger.NewLogger(cfg.Logger.Level, cfg.Logger.JSON)
	if err != nil {
		panic(err)
	}

	// Проверить режим запуска (HTTP API или CLI)
	if len(os.Args) > 1 && os.Args[1] != "serve" {
		runCLI(cfg, log, os.Args[1:])
		return
	}

	// Инициализировать БД
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", nil)
	}
	defer db.Close()

	// Создать репозитории
	messageRepo := postgres.NewMessageRepository(db)
	channelRepo := postgres.NewChannelRepository(db)
	digestRepo := postgres.NewDigestRepository(db)

	// Создать сервисы
	telegramSvc := telegram.NewTelegramService(log)

	// Создать обработчики команд
	generateDigestCmdHandler := commands.NewGenerateDigestCommandHandler(
		telegramSvc,
		messageRepo,
		channelRepo,
		digestRepo,
		log,
	)

	// Создать обработчики запросов
	getDigestQueryHandler := queries.NewGetDigestQueryHandler(digestRepo, log)

	// Создать HTTP обработчик
	handler := http.NewDigestHandler(generateDigestCmdHandler, getDigestQueryHandler, log)

	// Создать роутер
	router := http.NewRouter(handler, log)
	router.Setup()

	// Запустить HTTP сервер
	serverPort := strconv.Itoa(cfg.Server.Port)
	go func() {
		if err := router.Run(serverPort); err != nil {
			log.Error("server error", nil)
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("shutting down server", nil)
}

// runCLI запускает приложение в режиме CLI
func runCLI(cfg config.Config, log logger.Logger, args []string) {
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", nil)
	}
	defer db.Close()

	messageRepo := postgres.NewMessageRepository(db)
	channelRepo := postgres.NewChannelRepository(db)
	digestRepo := postgres.NewDigestRepository(db)

	telegramSvc := telegram.NewTelegramService(log)

	generateDigestCmdHandler := commands.NewGenerateDigestCommandHandler(
		telegramSvc,
		messageRepo,
		channelRepo,
		digestRepo,
		log,
	)

	cliApp := cli.NewCLIApp(generateDigestCmdHandler, log)
	if err := cliApp.Run(args); err != nil {
		log.Error(err.Error(), nil)
		os.Exit(1)
	}
}
```

## 3. Docker Configuration

### 3.1 Dockerfile

```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Установить необходимые инструменты
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Копировать go.mod и go.sum
COPY go.mod go.sum ./

# Загрузить зависимости
RUN go mod download

# Копировать исходный код
COPY . .

# Собрать приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram-digest ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

# Установить ca-certificates для HTTPS и tzdata для временных зон
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копировать бинарный файл из builder
COPY --from=builder /app/telegram-digest .

# Копировать миграции
COPY migrations/ ./migrations/

# Создать непривилегированного пользователя
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser
USER appuser

# Expose порт
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD /app/telegram-digest health || exit 1

# Запустить приложение
CMD ["/app/telegram-digest", "serve"]
```

### 3.2 docker-compose.yml

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    container_name: telegram-digest-db
    environment:
      POSTGRES_DB: ${DB_NAME:-telegram_digest}
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "${DB_PORT:-5432}:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: telegram-digest-app
    environment:
      TELEGRAM_API_ID: ${TELEGRAM_API_ID}
      TELEGRAM_API_HASH: ${TELEGRAM_API_HASH}
      TELEGRAM_PHONE: ${TELEGRAM_PHONE}
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER:-postgres}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME:-telegram_digest}
      SERVER_PORT: 8080
      LOG_LEVEL: ${LOG_LEVEL:-info}
      LOG_JSON: ${LOG_JSON:-false}
    ports:
      - "${SERVER_PORT:-8080}:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

## 4. Go Module Configuration

### go.mod

```
module telegram-news-digest

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.26.0
	github.com/gotd/td v0.87.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	// ... остальные зависимости
)
```

## 5. Configuration File

### config/config.yaml

```yaml
server:
  port: 8080
  read_timeout: 15
  write_timeout: 15
  idle_timeout: 60

database:
  host: localhost
  port: 5432
  user: postgres
  password: ${DB_PASSWORD}
  dbname: telegram_digest
  sslmode: disable

telegram:
  api_id: ${TELEGRAM_API_ID}
  api_hash: ${TELEGRAM_API_HASH}
  phone: ${TELEGRAM_PHONE}
  session_id: ${TELEGRAM_SESSION_ID}

logger:
  level: info
  json: false
```

### .env.example

```
# Telegram Configuration
TELEGRAM_API_ID=your_api_id
TELEGRAM_API_HASH=your_api_hash
TELEGRAM_PHONE=+1234567890
TELEGRAM_SESSION_ID=your_session_id

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=telegram_digest

# Server Configuration
SERVER_PORT=8080

# Logger Configuration
LOG_LEVEL=info
LOG_JSON=false
```

## 6. Makefile

```makefile
.PHONY: help build run test clean docker-build docker-up docker-down migrate lint fmt

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make migrate       - Run database migrations"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"

build:
	@echo "Building application..."
	CGO_ENABLED=0 go build -o bin/telegram-digest ./cmd/api

run: build
	@echo "Running application..."
	./bin/telegram-digest serve

test:
	@echo "Running tests..."
	go test -v -cover ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

docker-build:
	@echo "Building Docker image..."
	docker build -t telegram-digest:latest .

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

migrate:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/telegram_digest?sslmode=disable" up

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	gofmt -w .
	go mod tidy
```

## 7. Тестирование API

### Примеры curl команд

```bash
# Проверка здоровья приложения
curl http://localhost:8080/health

# Формирование дайджеста
curl -X POST http://localhost:8080/api/v1/digests/generate \
  -H "Content-Type: application/json" \
  -d '{
    "channel_username": "@mychannel",
    "from_date": "2024-01-01T00:00:00Z",
    "to_date": "2024-01-31T23:59:59Z"
  }'

# Получение дайджеста
curl http://localhost:8080/api/v1/digests/1
```

## 8. Запуск приложения

### Способ 1: Локальный запуск

```bash
# Установить зависимости
go mod download

# Установить tools
make fmt
make lint

# Запустить БД
docker run --name postgres -e POSTGRES_PASSWORD=password -d postgres:16-alpine

# Запустить приложение в режиме CLI
./bin/telegram-digest generate \
  -channel=@mychannel \
  -from=2024-01-01T00:00:00Z \
  -to=2024-01-31T23:59:59Z

# Или запустить HTTP сервер
./bin/telegram-digest serve
```

### Способ 2: Docker Compose

```bash
# Скопировать файл конфигурации
cp .env.example .env

# Отредактировать .env с вашими параметрами
# Запустить контейнеры
make docker-up

# Проверить логи
docker-compose logs -f app
```

## 9. Ключевые архитектурные принципы

### Clean Architecture
- **Entities**: Бизнес-логика сосредоточена в domain entities
- **Use Cases**: Application layer содержит команды и запросы (CQRS)
- **Interface Adapters**: HTTP и CLI адаптеры изолированы
- **Frameworks & Drivers**: Infrastructure слой для работы с внешними сервисами

### CQRS (Command Query Responsibility Segregation)
- **Commands**: GenerateDigestCommand для изменения состояния
- **Queries**: GetDigestQuery для чтения данных
- Разделение логики чтения и записи

### DDD (Domain-Driven Design)
- **Entities**: Message, Channel, Digest с инкапсулированной логикой
- **Value Objects**: DateRange для валидированных значений
- **Repositories**: Абстракция доступа к данным
- **Domain Services**: TelegramService для операций с Telegram API
- **Domain Errors**: Специфичные ошибки бизнес-логики

### Go Best Practices
- Использование интерфейсов для зависимостей
- Структурированное логирование с zap
- Обработка ошибок на каждом уровне
- Context для управления lifecycle операций
- Правильная организация кода по domain boundaries

