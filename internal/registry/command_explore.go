package registry

import (
	"fmt"

	"github.com/thalesraymond/pokedexcli/internal/api"
)

func commandExplore(cfg *PokedexContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore command requires a location area name")
	}
	locationName := args[0]
	client := api.NewPokedexClient()

	locationArea, err := client.GetLocationArea(locationName)

	if err != nil {
		return err
	}

	for _, pokemonEncounter := range locationArea.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}

	return nil
}
