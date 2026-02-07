package dto


type chatData struct {
	
}

type messageData struct {
	Id int64 `json:"message_id" binding:"required"`
	Chat chatData
}

type SaveWebhookBotUpdateRequest struct {
	UpdateId int64 `json:"update_id" binding:"required"`
	Message messageData `json:"message" binding:"required"`
}

