package postgres

import (
	"encoding/json"
	"strings"

	"telegram/internal/domain"
	"telegram/internal/domain/entities"
	"telegram/internal/domain/repositories"

	"github.com/jmoiron/sqlx"
)

type BotWebhookUpdateRepository struct {
	db *sqlx.DB
	repositories.BotWebhookUpdateRepositoryInterface
}

func NewBotWebhookUpdateRepository(db *sqlx.DB) repositories.BotWebhookUpdateRepositoryInterface {
	return &BotWebhookUpdateRepository{db: db}
}

func (rep *BotWebhookUpdateRepository) Save(ent *entities.BotWebhookUpdate) error {

	jsonData, errJson := json.Marshal(ent.Data)
	if errJson != nil {
		return errJson
	}

	_, err := rep.db.NamedExec(`INSERT INTO bot_webhooks_updates (update_id,bot_name,data,created_at) VALUES (:updateId,:botName,:data,:createdAt)`, 
        map[string]any{
            "updateId": ent.UpdateId,
            "botName": ent.BotName,
			"data": string(jsonData),
			"createdAt": ent.CreatedAt.Format("2006-01-02 15:04:05"),
    })

	if (err != nil) {

		if (strings.Contains(err.Error(), "duplicate key")) {
			return domain.ErrBotWebhookUpdateAlreadyExists
		}

	}

	return err
}