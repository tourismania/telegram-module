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
	"telegram/internal/presentation/http"
	"time"
)

func main() {
	fmt.Println("Погнали делать digest")

	message, err := entities.NewMessage(1, "dfdsfdsfsdf", time.Now())

	if err != nil {
		log.Fatal("Что-то какая-то хуита " + err.Error())
	}

	// Проверить режим запуска HTTP API
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		log.Println(os.Args[1])
		log.Println(message.GetText())

		saveWebhookCommandHandler := commands.NewSaveWebhookCommandHandler()

		handler := http.NewHandler(saveWebhookCommandHandler)

		router := http.NewRouter(*handler)
		router.Setup()

		// Запустить HTTP сервер
		serverPort := strconv.Itoa(8082)
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
