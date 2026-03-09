package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"telegram/internal/infrastructure/config"
	"telegram/internal/infrastructure/logger"
)

// NewConnection создает новое соединение с PostgreSQL
func NewConnection(cfg config.DatabaseConfig, logg logger.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logg.Error("failed to open database:" + err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logg.Error("failed to ping database:" + err.Error())
		return nil, err
	}

	return db, nil
}
