package commands

import (
	"telegram/internal/infrastructure/tg-bot-api"
)

// описываем объект команды
type SaveWebhookCommand struct {
	CommandName string
}

func NewSaveWebhookCommand(cn string) *SaveWebhookCommand {
	return &SaveWebhookCommand{
		CommandName: cn,
	}
}

// описываем обработчик команды
type SaveWebhookCommandHandler struct {
	tgBotApiService *tgbotapi.TgBotApiService
}

func NewSaveWebhookCommandHandler(tgBotApiService *tgbotapi.TgBotApiService)  *SaveWebhookCommandHandler {
	return &SaveWebhookCommandHandler{
		tgBotApiService: tgBotApiService,
	}
}

// метод обработчика команды
func (handler SaveWebhookCommandHandler) Handle(cmd SaveWebhookCommand) (bool) {
	return true
}

