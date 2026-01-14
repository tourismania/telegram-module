package entities

import (
	"time"
)

// Message представляет сообщение из Telegram канала
type Message struct {
	ID        int64
	ChannelID int64
	Text      string
	MediaURL  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMessage создает новое сообщение с валидацией
func NewMessage(channelID int64, text string, createdAt time.Time) (*Message, error) {

	if text == "" {
		//TODO: сделать ошибку 
		// return nil, ErrInvalidMessage
	}

	return &Message{
		ChannelID: channelID,
		Text:      text,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}, nil
}