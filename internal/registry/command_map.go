package registry

import (
	"fmt"
)

func commandMap(cfg *PokedexContext, args ...string) error {
	locationAreaResponse, err := cfg.Client.GetLocations(cfg.LocationAreasNextURL)

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
