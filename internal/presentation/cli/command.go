package cli

import (
	"fmt"
	"log"
)

type CliApp struct {
}

// создание cli приложения
func NewCliApp() *CliApp {
	return &CliApp{}
}

func (app CliApp) Run(args []string) error {

	if len(args) == 0 {
		app.printUsage()
		return fmt.Errorf("no command provided")
	}

	switch args[0] {
	case "didgest:generate-and-send":
		log.Printf("Digest generated successfully and sended",)
		return nil

	default:
		app.printUsage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}


// printUsage выводит справку по использованию
func (app CliApp) printUsage() {
	fmt.Println("Usage: telegram-module <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  generate  Generate a news digest for a Telegram channel")
	fmt.Println("\nExamples:")
	fmt.Println("  telegram-module generate -channel=@mychannel -from=2024-01-01T00:00:00Z -to=2024-01-31T23:59:59Z")
}