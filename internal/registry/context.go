package registry

import "github.com/thalesraymond/pokedexcli/internal/api"

type PokemonAPI interface {
	GetLocations(url *string) (api.LocationAreaResponse, error)
	GetLocationArea(name string) (api.LocationArea, error)
	GetPokemon(name string) (api.Pokemon, error)
}

type PokedexContext struct {
	Client                   PokemonAPI
	LocationAreasNextURL     *string
	LocationAreasPreviousURL *string
	Pokedex                  map[string]bool
}

func NewPokedexContext(apiClient PokemonAPI) *PokedexContext {
	return &PokedexContext{
		Client:                   apiClient,
		LocationAreasNextURL:     nil,
		LocationAreasPreviousURL: nil,
		Pokedex:                  make(map[string]bool),
	}
}
