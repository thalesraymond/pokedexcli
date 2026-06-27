package registry

import (
	"fmt"

	"github.com/thalesraymond/pokedexcli/api"
)

var previousUrl *string = nil

func commandMapBack() error {
	client := api.NewPokedexClient()

	locationAreaResponse, err := client.GetLocations(previousUrl)

	if err != nil {
		return err
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	if locationAreaResponse.Previous != nil {
		previousUrl = locationAreaResponse.Previous
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}
