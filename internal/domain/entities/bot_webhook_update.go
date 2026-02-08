package entities

import "time"

type BotWebhookUpdate struct {
	UpdateId int64
	BotName string
	Data any
	CreatedAt time.Time
}