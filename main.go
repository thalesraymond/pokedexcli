package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/thalesraymond/pokedexcli/internal/registry"
)

func main() {

	commands := registry.GetCLICommands()

	scanner := bufio.NewScanner(os.Stdin)

	pokedexContext := &registry.PokedexContext{}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()
		fields := strings.Fields(userInput)
		if len(fields) == 0 {
			continue
		}
		firstWord := strings.ToLower(fields[0])

		if command, ok := commands[firstWord]; ok {
			if err := command.Execute(pokedexContext, fields[1:]...); err != nil {
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
