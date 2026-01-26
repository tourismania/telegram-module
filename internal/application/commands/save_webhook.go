package commands

import (
	"fmt"
	"telegram/internal/infrastructure/tg-bot-api"
)

type SaveWebhookCommand struct {
	CommandName string
}

func NewSaveWebhookCommand(cn string) *SaveWebhookCommand {
	return &SaveWebhookCommand{
		CommandName: cn,
	}
}

type SaveWebhookCommandHandler struct {
	tgBotApiService *tgbotapi.TgBotApiService
}

func NewSaveWebhookCommandHandler(tgBotApiService *tgbotapi.TgBotApiService)  *SaveWebhookCommandHandler {
	return &SaveWebhookCommandHandler{
		tgBotApiService: tgBotApiService,
	}
}

func (handler SaveWebhookCommandHandler) Handle(cmd SaveWebhookCommand) (bool) {
	fmt.Println("Ура, добрались до команды!!!!!");
	fmt.Println(cmd);

	handler.tgBotApiService.PrintUpdates()
	return true
}

