package repositories

import (
	"telegram/internal/domain/entities"
)

type BotWebhookUpdateRepositoryInterface interface {
	Save(*entities.BotWebhookUpdate) error
}