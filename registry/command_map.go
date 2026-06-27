package registry

import (
	"fmt"

	"github.com/thalesraymond/pokedexcli/api"
)

var nextURL *string = nil

func commandMap() error {
	client := api.NewPokedexClient()

	locationAreaResponse, err := client.GetLocations(nextURL)

	if err != nil {
		return err
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	if locationAreaResponse.Next != nil {
		nextURL = locationAreaResponse.Next
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}
