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