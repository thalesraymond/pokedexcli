package api

import (
	"encoding/json"
	"io"
)

type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             NamedAPIResource      `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource         `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetail `json:"version_details"`
}

type EncounterVersionDetail struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource          `json:"pokemon"`
	VersionDetails []PokemonEncounterVersion `json:"version_details"`
}

type PokemonEncounterVersion struct {
	Version          NamedAPIResource  `json:"version"`
	MaxChance        int               `json:"max_chance"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type EncounterDetail struct {
	MinLevel        int                `json:"min_level"`
	MaxLevel        int                `json:"max_level"`
	ConditionValues []NamedAPIResource `json:"condition_values"`
	Chance          int                `json:"chance"`
	Method          NamedAPIResource   `json:"method"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func (c *PokedexClient) GetLocations(url *string) (LocationAreaResponse, error) {

	urlStr := Endpoints["location-area"]
	if url != nil {
		urlStr = *url
	}

	if val, ok := c.cache.Get(urlStr); ok {
		var locationAreaResponse LocationAreaResponse
		err := json.Unmarshal(val, &locationAreaResponse)
		if err != nil {
			return LocationAreaResponse{}, err
		}
		return locationAreaResponse, nil
	}

	res, err := c.httpClient.Get(urlStr)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	c.cache.Set(urlStr, body)

	var locationAreaResponse LocationAreaResponse
	err = json.Unmarshal(body, &locationAreaResponse)

	if err != nil {
		return locationAreaResponse, err
	}

	return locationAreaResponse, nil
}

func (c *PokedexClient) GetLocationArea(locationAreaName string) (LocationArea, error) {
	urlStr := Endpoints["location-area"] + "/" + locationAreaName

	if val, ok := c.cache.Get(urlStr); ok {
		var locationAreaDto LocationArea
		err := json.Unmarshal(val, &locationAreaDto)
		if err != nil {
			return LocationArea{}, err
		}
		return locationAreaDto, nil
	}

	res, err := c.httpClient.Get(urlStr)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}

	c.cache.Set(urlStr, body)

	var locationAreaDto LocationArea
	err = json.Unmarshal(body, &locationAreaDto)

	if err != nil {
		return locationAreaDto, err
	}

	return locationAreaDto, nil
}
