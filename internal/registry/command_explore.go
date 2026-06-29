package registry

import (
	"fmt"
)

func commandExplore(context *PokedexContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore command requires a location area name")
	}
	locationName := args[0]

	locationArea, err := context.Client.GetLocationArea(locationName)

	if err != nil {
		return err
	}

	for _, pokemonEncounter := range locationArea.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}

	return nil
}
