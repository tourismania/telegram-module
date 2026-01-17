package queries

type GetMessagesQuery struct {}

type GetMessagesQueryHandler struct {}

func NewGetMessagesQueryHandler() *GetMessagesQueryHandler {
	return &GetMessagesQueryHandler{}
}

// TODO: нужно возвращать слайс из сообщений
func (h GetMessagesQueryHandler) Handle() {



}