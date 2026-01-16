package main

import (
	"fmt"
	"log"
	"os"
	"telegram/internal/domain/entities"
	"time"
	"telegram/internal/presentation/cli"
)

func main() {
	fmt.Println("Погнали делать digest")

	message, err := entities.NewMessage(1, "dfdsfdsfsdf", time.Now())

	if err != nil {
		log.Fatal("Что-то какая-то хуита " + err.Error())
	}

	if len(os.Args) <= 1 {
		log.Fatal("Не указаны параметры запуска")
	}

	// Проверить режим запуска (HTTP API или CLI)
	cliApp := cli.NewCliApp()
	if err := cliApp.Run(os.Args[1:]); err != nil {
		log.Fatalln(err.Error(), nil)
		os.Exit(1)
	}
	log.Println(message.GetText())

}
