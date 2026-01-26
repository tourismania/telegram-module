package tgbotapi

import (
	"log"

	tgbotapiv5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotApiService struct {
	token string
}

func NewTgBotApiService(t string) *TgBotApiService {
	return &TgBotApiService{
		token: t,
	};
}

func (srvc *TgBotApiService) PrintUpdates() {
	bot, err := tgbotapiv5.NewBotAPI(srvc.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapiv5.NewUpdate(1)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}
	}
}