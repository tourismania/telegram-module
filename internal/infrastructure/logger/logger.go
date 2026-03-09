// Package logger предоставляет абстракцию логирования на базе zap с настраиваемым уровнем.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger — интерфейс логирования с поддержкой уровней Debug, Info, Warn, Error и Fatal.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// zapLogger — реализация Logger на базе zap.
type zapLogger struct {
	log *zap.Logger
}

// NewLogger создаёт логгер с заданным уровнем (debug, info, warn, error).
// Используется production-конфигурация zap.
func NewLogger(level string) (Logger, error) {
	config := zap.NewProductionConfig()

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
	defer log.Sync()

	return &zapLogger{log: log}, nil
}

// Debug записывает отладочное сообщение.
func (zl *zapLogger) Debug(msg string) {
	zl.log.Debug(msg)
}

// Info записывает информационное сообщение.
func (zl *zapLogger) Info(msg string) {
	zl.log.Info(msg)
}

// Warn записывает предупреждение.
func (zl *zapLogger) Warn(msg string) {
	zl.log.Warn(msg)
}

// Error записывает сообщение об ошибке.
func (zl *zapLogger) Error(msg string) {
	zl.log.Error(msg)
}

// Fatal записывает критическое сообщение и завершает программу (os.Exit(1)).
func (zl *zapLogger) Fatal(msg string) {
	zl.log.Fatal(msg)
}