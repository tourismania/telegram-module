package postgres

import (
	"strings"

	"telegram/internal/domain"
	"telegram/internal/domain/entities"
	"telegram/internal/domain/repositories"

	"github.com/jmoiron/sqlx"
)

type BotWebhookUpdateRepository struct {
	db *sqlx.DB
}
// проверит соответствие интерфейсу в compile-time
var _ repositories.BotWebhookUpdateRepositoryInterface = (*BotWebhookUpdateRepository)(nil)

func NewBotWebhookUpdateRepository(db *sqlx.DB) repositories.BotWebhookUpdateRepositoryInterface {
	return &BotWebhookUpdateRepository{db: db}
}

func (r *BotWebhookUpdateRepository) Save(ent entities.BotWebhookUpdate) error {

	_, err := r.db.NamedExec(`INSERT INTO bot_webhooks_updates (update_id,bot_name,payload,created_at) VALUES (:updateId,:botName,:payload,:createdAt)`, 
        map[string]any{
            "updateId": ent.UpdateId(),
            "botName": ent.BotName(),
			"payload": ent.Payload(),
			"createdAt": ent.CreatedAt().Format("2006-01-02 15:04:05"),
    })

	if (err != nil) {
		if (strings.Contains(err.Error(), "duplicate key")) {
			return domain.ErrBotWebhookUpdateAlreadyExists
		}
	}

	return err
}