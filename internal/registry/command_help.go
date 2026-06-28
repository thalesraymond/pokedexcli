package registry

import (
	"fmt"
)

func commandHelp(cfg *PokedexContext) error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	commands := GetCLICommands()
	for _, command := range commands {
		fmt.Printf("  %s: %s\n", command.name, command.description)
	}
	return nil
}
