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

	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		firstWord := strings.TrimSpace(strings.ToLower(strings.Fields(userInput)[0]))

		if command, ok := commands[firstWord]; ok {
			if err := command.Execute(); err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
