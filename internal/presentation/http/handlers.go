package http

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
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
func (h *Handler) SaveWebhookBotUpdate(c *gin.Context) {
	var req dto.SaveWebhookBotUpdateRequest

	// Читаем тело с лимитом 1MB, так как необходимо сохранять всю информацию
    bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }
    defer c.Request.Body.Close()

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
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

		// Тело как строка (или []byte)
    rawBody := string(bodyBytes)
	rawBody = strings.ReplaceAll(rawBody, " ", "")  // Удалит пробелы
	rawBody = strings.ReplaceAll(rawBody, "\n", "") // Удалит переносы
	rawBody = strings.ReplaceAll(rawBody, "\t", "") // Удалит табуляции

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
		Payload: rawBody,
	}

	res, err := h.saveBotWebhookUpdate.Handle(&command)

	if (err != nil) {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &dto.SaveWebhookBotUpdateResponse{Status: res})
}
