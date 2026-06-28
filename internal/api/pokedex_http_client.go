package api

import (
	"net/http"
	"time"
)

// PokedexClient is a struct that holds the state of the crawler client
type PokedexClient struct {
	userAgent  string
	timeout    time.Duration
	httpClient *http.Client
}

// NewPokedexClient creates a new instance of PokedexClient
func NewPokedexClient() *PokedexClient {
	timeout := 10 * time.Second
	return &PokedexClient{
		userAgent: "pokedex-cli/1.0",
		timeout:   timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}
