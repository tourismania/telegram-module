package dto

// SaveWebhook DTO для запроса на формирование дайджеста
type SaveWebhooRequest struct {
	CommandName string    `json:"command" binding:"required"`
}