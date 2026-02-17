package entities

import "time"

type BotWebhookUpdate struct {
	updateId int64
	botName string
	payload any
	createAt time.Time
}

func (ent BotWebhookUpdate) UpdateId() int64 {
	return ent.updateId
}

func (ent BotWebhookUpdate) BotName() string {
	return ent.botName
}

func (ent BotWebhookUpdate) Payload() any {
	return ent.payload
}

func (ent BotWebhookUpdate) CreatedAt() time.Time {
	return ent.createAt
}

func NewBookWebhookUpdate(updateId int64, botName string, payload any, createdAt time.Time) *BotWebhookUpdate {
	return &BotWebhookUpdate{
		updateId: updateId,
		botName: botName,
		payload: payload,
		createAt: createdAt,
	}
}