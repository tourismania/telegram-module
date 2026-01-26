package dto

// SaveWebhookResponse DTO для ответа с дайджестом
type SaveWebhookResponse struct {
	Status bool
}

type ErrorResponse struct {
	Error string
	Message string
}