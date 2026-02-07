package commands

import (
	"fmt"
	"telegram/internal/infrastructure/tg-bot-api"
)

// описываем объект команды
type SaveWebhookBotUpdateCommand struct {
	UpdateId int64
	Message messageInfo 
}

type messageInfo struct {
	Date int64
	Chat chatInfo
	MessageId int64
	From fromInfo
	Text string
}

type chatInfo struct {
	LastName string
	Id int64
	Type string
	FirstName string
	UserName string
}

type fromInfo struct {
	LastName string
	Id int64
	FirstName string
	UserName string
}

func NewSaveWebhookBotUpdateCommand(updateId int64, messageId int64) *SaveWebhookBotUpdateCommand {
	return &SaveWebhookBotUpdateCommand{
		UpdateId: int64(updateId),
		Message: messageInfo{
			Date: int64(100),
			Chat: chatInfo{
				LastName: "string",
				Id: int64(100),
				Type: "string",
				FirstName: "string",
				UserName: "string",
			},
			MessageId: messageId,
			From: fromInfo{
				LastName: "string",
				Id: int64(100),
				FirstName: "string",
				UserName: "string",
			},
			Text: "string",
		},
	}
}

// описываем обработчик команды
type SaveWebhookBotUpdateCommandHandler struct {
	tgBotApiService *tgbotapi.TgBotApiService
}

func NewSaveWebhookBotUpdateCommandHandler(tgBotApiService *tgbotapi.TgBotApiService)  *SaveWebhookBotUpdateCommandHandler {
	return &SaveWebhookBotUpdateCommandHandler{
		tgBotApiService: tgBotApiService,
	}
}

// метод обработчика команды
func (handler SaveWebhookBotUpdateCommandHandler) Handle(cmd SaveWebhookBotUpdateCommand) (bool) {
	
	fmt.Println(cmd.UpdateId)
	fmt.Println(cmd.Message.MessageId)
	return true
}

