package save_webhook_bot_update

import (
	"log"
	"telegram/internal/domain/entities"
	"telegram/internal/domain/repositories"
	"time"
)

// описываем обработчик команды
type Handler struct {
	webhookBotUpdateRepository repositories.BotWebhookUpdateRepositoryInterface
}

func NewHandler(rep repositories.BotWebhookUpdateRepositoryInterface) *Handler {
	return &Handler{
		webhookBotUpdateRepository: rep,
	}
}

// метод обработчика команды
func (h *Handler) Handle(cmd Command) (error) {

	ent := entities.NewBookWebhookUpdate(
		cmd.UpdateId,
		cmd.BotName,
		cmd.Payload,
		time.Now(),
	)

	err := h.webhookBotUpdateRepository.Save(ent)

	if (err != nil) {
		log.Println("Error: " + err.Error())
	}
	
	return err
}

