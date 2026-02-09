package domain

import "errors"

var (
	ErrMessageInvalidText = errors.New("Invalid Text Message")
	ErrBotWebhookUpdateAlreadyExists = errors.New("Bot Webhook Update Already Exists!")
)