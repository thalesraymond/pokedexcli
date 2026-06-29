package registry

import (
	"fmt"
)

func commandPokedex(cfg *PokedexContext, args ...string) error {
	for _, pokemon := range cfg.Pokedex {
		fmt.Println(pokemon.Name)
	}

	return nil
}
