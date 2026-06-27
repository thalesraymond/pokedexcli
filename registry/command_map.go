package registry

import (
	"fmt"

	"github.com/thalesraymond/pokedexcli/api"
)

func commandMap(cfg *PokedexContext) error {
	client := api.NewPokedexClient()

	locationAreaResponse, err := client.GetLocations(cfg.LocationAreasNextURL)

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
