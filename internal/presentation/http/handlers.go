package http

import (
	"net/http"
	"os"
	"telegram/internal/application/commands"
	"telegram/internal/presentation/http/dto"

	"github.com/gin-gonic/gin"
)

const (
	headerTelegramBotApiSecretToken = "X-Telegram-Bot-Api-Secret-Token"
	envTelegramBotApiSecretToken = "TELEGRAM_BOT_API_SECRET_TOKEN"
)

// правильнее называть именно Handler, а не контроллер, так как так принято в Go
type Handler struct {
	saveBotWebhookUpdate *commands.SaveBotWebhookUpdateCommandHandler
}

func NewHandler(
	saveWebhookBotUpdateCommandHandler *commands.SaveBotWebhookUpdateCommandHandler,
) *Handler {
	return &Handler{
		saveBotWebhookUpdate: saveWebhookBotUpdateCommandHandler,
	}
}


// сохранение вэбхуки из телеграм бота
func (h *Handler) SaveWebhook(c *gin.Context) {
	var req dto.SaveWebhookBotUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	botName := "undefined";
	if (os.Getenv(envTelegramBotApiSecretToken) == c.GetHeader(headerTelegramBotApiSecretToken)) {
		botName = "tourismania"
	}

	// создаим команду
	command := commands.SaveBotWebhookUpdateCommand{
		UpdateId: req.UpdateId,
		Message: commands.Message{
			Date: req.Message.Date,
			MessageId: req.Message.MessageId,
			Text: req.Message.Text,
			Chat: commands.Chat{
				LastName: req.Message.Chat.LastName,
				FirstName: req.Message.Chat.FirstName,
				Id: req.Message.Chat.Id,
				Type: req.Message.Chat.Type,
				Username: req.Message.Chat.Username,
			},
			From: commands.From{
				LastName: req.Message.From.LastName,
				FirstName: req.Message.From.FirstName,
				Id: req.Message.From.Id,
				Username: req.Message.From.Username,
			},
		},
		BotName: botName,
	}

	// выполним команду
	r := h.saveBotWebhookUpdate.Handle(&command)

	c.JSON(http.StatusOK, &dto.SaveWebhookBotUpdateResponse{Status: r})
}
