package dto

type from struct {
	Id int64 `json:"id"`
	LastName string `json:"last_name"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
}

type chat struct {
	LastName string `json:"last_name"`
	Id int64 `json:"id"`
	Type string `json:"type"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
}

type message struct {
	MessageId int64 `json:"message_id"`
	Chat chat `json:"chat"`
	Date int64 `json:"date"`
	From from `json:"from"`
	Text string `json:"text"`
}

type SaveWebhookBotUpdateRequest struct {
	UpdateId int64 `json:"update_id" binding:"required"`
	Message message `json:"message"`
}

