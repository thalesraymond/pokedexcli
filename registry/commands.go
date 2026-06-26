package registry

import (
	"fmt"
	"os"
)

func GetCLICommands() map[string]CLICommand {
	commands := map[string]CLICommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
	}

	return commands
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	commands := GetCLICommands()
	for _, command := range commands {
		fmt.Printf("  %s: %s\n", command.name, command.description)
	}
	return nil
}
