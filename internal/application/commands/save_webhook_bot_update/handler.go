package save_webhook_bot_update

import (
	"telegram/internal/domain/entities"
	"telegram/internal/domain/repositories"
	"telegram/internal/infrastructure/logger"
	"time"
)

// описываем обработчик команды
type Handler struct {
	webhookBotUpdateRepository repositories.BotWebhookUpdateRepositoryInterface
	logg logger.Logger
}

func NewHandler(rep repositories.BotWebhookUpdateRepositoryInterface, logg logger.Logger) *Handler {
	return &Handler{
		webhookBotUpdateRepository: rep,
		logg: logg,
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

	if err != nil {
		h.logg.Error(err.Error())
	}
	
	return err
}

