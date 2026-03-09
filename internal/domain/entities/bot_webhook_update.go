package entities

import "time"

type BotWebhookUpdate struct {
	updateId int64
	botName string
	payload any
	createAt time.Time
}

func (b BotWebhookUpdate) UpdateId() int64 {
	return b.updateId
}

func (b BotWebhookUpdate) BotName() string {
	return b.botName
}

func (b BotWebhookUpdate) Payload() any {
	return b.payload
}

func (b BotWebhookUpdate) CreatedAt() time.Time {
	return b.createAt
}

func NewBookWebhookUpdate(updateId int64, botName string, payload any, createdAt time.Time) BotWebhookUpdate {
	return BotWebhookUpdate{
		updateId: updateId,
		botName: botName,
		payload: payload,
		createAt: createdAt,
	}
}