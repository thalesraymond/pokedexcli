package registry

import (
	"fmt"
)

// Name: pidgey
// Height: 3
// Weight: 18
// Stats:
//   -hp: 40
//   -attack: 45
//   -defense: 40
//   -special-attack: 35
//   -special-defense: 35
//   -speed: 56
// Types:
//   - normal
//   - flying

func commandInspect(cfg *PokedexContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("inspect command requires a pokemon name")
	}
	pokemonName := args[0]

	if pokemon, ok := cfg.Pokedex[pokemonName]; ok {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)

		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		for _, pokemonType := range pokemon.Types {
			fmt.Printf("  -%s\n", pokemonType.Type.Name)
		}

		return nil
	}

	fmt.Printf("Pokemon %s not found in Pokedex\n", pokemonName)

	return nil
}
