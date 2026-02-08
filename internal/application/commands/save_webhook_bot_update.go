package commands

import (
	"log"
	"telegram/internal/domain/entities"
	"telegram/internal/domain/repositories"
	"time"
)

// описываем объект команды
type SaveBotWebhookUpdateCommand struct {
	UpdateId int64
	Message Message 
	BotName string
}

type Message struct {
	Date int64
	Chat Chat
	MessageId int64
	From From
	Text string
}

type Chat struct {
	LastName string
	Id int64
	Type string
	FirstName string
	Username string
}

type From struct {
	LastName string
	Id int64
	FirstName string
	Username string
}

// описываем обработчик команды
type SaveBotWebhookUpdateCommandHandler struct {
	webhookBotUpdateRepository repositories.BotWebhookUpdateRepositoryInterface
}

func NewSaveBotWebhookUpdateCommandHandler(bwuri repositories.BotWebhookUpdateRepositoryInterface)  *SaveBotWebhookUpdateCommandHandler {
	return &SaveBotWebhookUpdateCommandHandler{
		webhookBotUpdateRepository: bwuri,
	}
}

// метод обработчика команды
func (handler SaveBotWebhookUpdateCommandHandler) Handle(cmd *SaveBotWebhookUpdateCommand) (bool) {

	ent := entities.BotWebhookUpdate{
		UpdateId: cmd.UpdateId,
		BotName: cmd.BotName,
		Data: cmd.Message,
		CreatedAt: time.Now(),
	}

	err := handler.webhookBotUpdateRepository.Save(&ent)

	if (err != nil) {
		log.Fatalln("ALARM ERROR")
		log.Fatalln(err.Error())
	}
	
	return true
}

