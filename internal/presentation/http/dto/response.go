package dto

type SaveWebhookBotUpdateResponse struct {
	Status bool `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
}