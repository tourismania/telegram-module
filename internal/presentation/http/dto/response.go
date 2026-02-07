package dto

type SaveWebhookBotUpdateResponse struct {
	Status bool `json:"status"`
}

type ErrorResponse struct {
	Error string
	Message string
}