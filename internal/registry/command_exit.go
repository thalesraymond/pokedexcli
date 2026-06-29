package registry

import (
	"fmt"
	"os"
)

func commandExit(context *PokedexContext, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)	
	return nil
}
