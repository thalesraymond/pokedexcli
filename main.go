package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/thalesraymond/pokedexcli/registry"
)

func main() {

	commands := registry.GetCLICommands()

	scanner := bufio.NewScanner(os.Stdin)

	pokedexContext := &registry.PokedexContext{}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		firstWord := strings.TrimSpace(strings.ToLower(strings.Fields(userInput)[0]))

		if command, ok := commands[firstWord]; ok {
			if err := command.Execute(pokedexContext); err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
