package registry

import (
	"fmt"
	"os"
)

func commandExit(cfg *PokedexContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)	
	return nil
}
