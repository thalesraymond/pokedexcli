package api

import (
	"net/http"
	"time"

	"github.com/thalesraymond/pokedexcli/internal/cache"
)

// PokedexClient is a struct that holds the state of the crawler client
type PokedexClient struct {
	userAgent  string
	timeout    time.Duration
	httpClient *http.Client
	cache      *cache.CacheData
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
		cache: cache.NewCacheData(),
	}
}
