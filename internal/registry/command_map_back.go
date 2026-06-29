package registry

import (
	"fmt"
)

func commandMapBack(context *PokedexContext, args ...string) error {
	if context.LocationAreasPreviousURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreaResponse, err := context.Client.GetLocations(context.LocationAreasPreviousURL)

	if err != nil {
		return err
	}

	for _, locationArea := range locationAreaResponse.Results {
		fmt.Println(locationArea.Name)
	}

	if locationAreaResponse.Next != nil {
		context.LocationAreasNextURL = locationAreaResponse.Next
	}
	if locationAreaResponse.Previous != nil {
		context.LocationAreasPreviousURL = locationAreaResponse.Previous
	}



	return nil
}
