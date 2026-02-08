package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	Telegram TelegramConfig
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type TelegramConfig struct {
	BotApiToken   string
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

	// Load the .env file
	// The function will look for a file named ".env" in the current directory by default.
	err := godotenv.Load()
	if (err != nil) {
		log.Print("Undefined env file")
	}

	return Config{
		Telegram: TelegramConfig{
			BotApiToken:  os.Getenv("TELEGRAM_BOT_API_TOKEN"),
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
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSL_MODE", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
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