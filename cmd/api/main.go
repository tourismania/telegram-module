package main

import (
	"os"

	"os/signal"
	"strconv"
	"syscall"
	"telegram/internal/application/commands/save_webhook_bot_update"
	"telegram/internal/infrastructure/config"
	"telegram/internal/infrastructure/logger"
	"telegram/internal/infrastructure/storage/postgres"
	"telegram/internal/presentation/http"
)

func main() {

	// подгрузим конфигурацию
	cfg := config.LoadConfig()

	// инициалиризуем logger
	logg, err := logger.NewLogger(cfg.Logger.LogLevel)
	if err != nil {
		logg.Fatal(err.Error())
	}

	// инициализируем коннект к БД
	db, err := postgres.NewConnection(cfg.Database, logg)
	if err != nil {
		logg.Fatal("failed to connect to database: " + err.Error())
	}
	defer db.Close()
	
	// инициализируем сервисы
	// tgBotApiService := tgbotapi.NewTgBotApiService(cnfg.Telegram.BotApiToken)

	// инициализируем репозитории
	botWebhookUpdateRepository := postgres.NewBotWebhookUpdateRepository(db)

	// инициализируем commands (Cqrs)
	saveWebhookBotUpdateCommandHandler := save_webhook_bot_update.NewHandler(botWebhookUpdateRepository, logg)

	// инициализируем обработчик для роутера
	handler := http.NewHandler(
		saveWebhookBotUpdateCommandHandler,
		cfg.Telegram,
	)

	// инициализируем роутер
	router := http.NewRouter(*handler)
	router.Setup()

	// Запустить HTTP сервер
	serverPort := strconv.Itoa(cfg.Server.Port)
	go func() {
		if err := router.Run(serverPort); err != nil {
			logg.Fatal("Error run server: " + err.Error())
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logg.Info("shutting down server")
}
