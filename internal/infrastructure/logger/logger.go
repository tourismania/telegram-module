package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger интерфейс логирования
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Field struct {
	Key       string
    Type      uint8
    Integer   int64
    String    string
    Interface any
}

type zapLogger struct {
	log *zap.Logger
}

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

func (zl *zapLogger) Debug(msg string) {
	zl.log.Debug(msg)
}

func (zl *zapLogger) Info(msg string) {
	zl.log.Info(msg)
}

func (zl *zapLogger) Warn(msg string) {
	zl.log.Warn(msg)
}

func (zl *zapLogger) Error(msg string) {
	zl.log.Error(msg)
}

func (zl *zapLogger) Fatal(msg string) {
	zl.log.Fatal(msg)
}