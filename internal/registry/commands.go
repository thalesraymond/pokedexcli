package registry

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

		"map": {
			name:        "map",
			description: "Display the available locations",
			callback:    commandMap,
		},

		"mapback": {
			name:        "mapback",
			description: "Display the previous available locations",
			callback:    commandMapBack,
		},

		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon",
			callback:    commandPokedex,
		},
	}

	return commands
}
