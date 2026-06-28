package registry

import (
	"fmt"

	"github.com/thalesraymond/pokedexcli/internal/api"
)

func commandMapBack(cfg *PokedexContext, args ...string) error {
	client := api.NewPokedexClient()

	if cfg.LocationAreasPreviousURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreaResponse, err := client.GetLocations(cfg.LocationAreasPreviousURL)

	if err != nil {
		return err
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	if locationAreaResponse.Next != nil {
		cfg.LocationAreasNextURL = locationAreaResponse.Next
	}
	if locationAreaResponse.Previous != nil {
		cfg.LocationAreasPreviousURL = locationAreaResponse.Previous
	}



	return nil
}
