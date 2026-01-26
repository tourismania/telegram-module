package http

import (
	"net/http"
	"telegram/internal/application/commands"
	"telegram/internal/application/dto"

	"github.com/gin-gonic/gin"
)

// правильнее называть именно Handler, а не контроллер, так как так принято в Go
type Handler struct {
	saveWebhook *commands.SaveWebhookCommandHandler
}

func NewHandler(
	saveWebhookHandler *commands.SaveWebhookCommandHandler,
) *Handler {
	return &Handler{
		saveWebhook: saveWebhookHandler,
	}
}

func (h *Handler) SaveWebhook(c *gin.Context) {
	var req dto.SaveWebhooRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Выполнить команду
	h.saveWebhook.Handle(*commands.NewSaveWebhookCommand(req.CommandName))

}
