package entities

import (
	"telegram/internal/domain"
	"time"
)

type Message struct {
	id int64
	text string
	createdAt time.Time
}

// создаем новое сообщение
func NewMessage(id int64, text string, createdAt time.Time) (*Message, error) {

	if text == "" {
		return nil, domain.ErrMessageInvalidText
	}

	return &Message{
		id: id,
		text: text,
		createdAt: time.Now(),
	}, nil
}

func (m Message) GetText () string {
	return m.text
}