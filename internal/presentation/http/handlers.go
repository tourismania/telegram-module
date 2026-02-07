package http

import (
	"net/http"
	"telegram/internal/application/commands"
	"telegram/internal/presentation/http/dto"
	"github.com/gin-gonic/gin"
)

// правильнее называть именно Handler, а не контроллер, так как так принято в Go
type Handler struct {
	saveWebhookBotUpdate *commands.SaveWebhookBotUpdateCommandHandler
}

func NewHandler(
	saveWebhookBotUpdateCommandHandler *commands.SaveWebhookBotUpdateCommandHandler,
) *Handler {
	return &Handler{
		saveWebhookBotUpdate: saveWebhookBotUpdateCommandHandler,
	}
}


func (h *Handler) SaveWebhook(c *gin.Context) {
	var req dto.SaveWebhookBotUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Выполнить команду
	r := h.saveWebhookBotUpdate.Handle(*commands.NewSaveWebhookBotUpdateCommand(req.UpdateId, req.Message.Id))

	c.JSON(http.StatusOK, &dto.SaveWebhookBotUpdateResponse{Status: r})
}
