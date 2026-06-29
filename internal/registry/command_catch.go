package registry

import (
	"fmt"
	"math/rand"

	"github.com/thalesraymond/pokedexcli/internal/api"
)

func commandCatch(cfg *PokedexContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch command requires a pokemon name")
	}
	pokemonName := args[0]
	client := api.NewPokedexClient()

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := client.GetPokemon(pokemonName)

	if err != nil {
		return err
	}

	randomNumber := rand.Intn(pokemon.BaseExperience)

	if randomNumber < 50 {
		fmt.Printf("You caught %s!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}
