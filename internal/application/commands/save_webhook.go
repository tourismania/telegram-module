package commands

import "fmt"

type SaveWebhookCommand struct {
	CommandName string
}

func NewSaveWebhookCommand(cn string) *SaveWebhookCommand {
	return &SaveWebhookCommand{
		CommandName: cn,
	}
}

type SaveWebhookCommandHandler struct {}

func NewSaveWebhookCommandHandler()  *SaveWebhookCommandHandler {
	return &SaveWebhookCommandHandler{}
}

func (handler SaveWebhookCommandHandler) Handle(cmd SaveWebhookCommand) (bool) {
	fmt.Println("Ура, добрались до команды!!!!!");
	fmt.Println(cmd);
	return true
}

