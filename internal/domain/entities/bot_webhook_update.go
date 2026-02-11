package entities

import "time"

type BotWebhookUpdate struct {
	UpdateId int64
	BotName string
	Payload any
	CreatedAt time.Time
}