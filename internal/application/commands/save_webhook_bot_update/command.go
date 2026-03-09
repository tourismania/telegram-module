package save_webhook_bot_update

// описываем объект команды
type Command struct {
	UpdateId int64
	Message Message 
	BotName string
	Payload any
}

type Message struct {
	Date int64
	Chat Chat
	MessageId int64
	From From
	Text string
}

type Chat struct {
	LastName string
	Id int64
	Type string
	FirstName string
	Username string
}

type From struct {
	LastName string
	Id int64
	FirstName string
	Username string
}