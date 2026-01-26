package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"telegram/internal/application/commands"
	"telegram/internal/domain/entities"
	"telegram/internal/infrastructure/tg-bot-api"
	"telegram/internal/infrastructure/config"
	"telegram/internal/presentation/http"
	"time"
)

func main() {
	fmt.Println("Погнали делать digest")

	message, err := entities.NewMessage(1, "dfdsfdsfsdf", time.Now())

	if err != nil {
		log.Fatal("Что-то какая-то хуита " + err.Error())
	}

	if len(os.Args) > 1 && os.Args[1] == "test" { 
		
	}

	// Проверить режим запуска HTTP API
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		log.Println(os.Args[1])
		log.Println(message.GetText())

		cnfg := config.LoadConfig()
		
		tgBotApiService := tgbotapi.NewTgBotApiService(cnfg.Telegram.BotApiToken)

		saveWebhookCommandHandler := commands.NewSaveWebhookCommandHandler(tgBotApiService)

		handler := http.NewHandler(saveWebhookCommandHandler)

		router := http.NewRouter(*handler)
		router.Setup()

		// Запустить HTTP сервер
		serverPort := strconv.Itoa(8088)
		go func() {
			if err := router.Run(serverPort); err != nil {
				log.Fatal("FATAALLLLL")
			}
		}()

		// Обработка сигналов завершения
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("shutting down server", nil)
		
		return
	}

	log.Fatal("Приложение работает только в режиме запуска сервера")

}
