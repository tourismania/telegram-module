package main

import (
	"log"
	"os"

	"os/signal"
	"strconv"
	"syscall"
	"telegram/internal/application/commands/save_webhook_bot_update"
	"telegram/internal/infrastructure/config"
	"telegram/internal/infrastructure/storage/postgres"
	"telegram/internal/presentation/http"
)

func main() {

	// подгрузим конфигурацию
	cfg := config.LoadConfig()

	// инициализируем коннект к БД
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", nil)
	}
	defer db.Close()
	
	// инициализируем сервисы
	// tgBotApiService := tgbotapi.NewTgBotApiService(cnfg.Telegram.BotApiToken)

	// инициализируем репозитории
	botWebhookUpdateRepository := postgres.NewBotWebhookUpdateRepository(db)

	// инициализируем commands (Cqrs)
	saveWebhookBotUpdateCommandHandler := save_webhook_bot_update.NewHandler(botWebhookUpdateRepository)

	// инициализируем обработчик для роутера
	handler := http.NewHandler(saveWebhookBotUpdateCommandHandler)

	// инициализируем роутер
	router := http.NewRouter(*handler)
	router.Setup()

	// Запустить HTTP сервер
	serverPort := strconv.Itoa(cfg.Server.Port)
	go func() {
		if err := router.Run(serverPort); err != nil {
			log.Fatal("Error run server")
			log.Fatalln(err)
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down server", nil)
	
}
