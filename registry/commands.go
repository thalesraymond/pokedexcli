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
	}

	return commands
}
