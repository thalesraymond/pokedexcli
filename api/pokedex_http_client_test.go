package api

import (
	"testing"
	"time"
)

// --- NewPokedexClient ---

func TestNewPokedexClient_Defaults(t *testing.T) {
	client := NewPokedexClient()

	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.userAgent != "pokedex-cli/1.0" {
		t.Errorf("unexpected userAgent: got %q, want %q", client.userAgent, "pokedex-cli/1.0")
	}
	if client.timeout != 10*time.Second {
		t.Errorf("unexpected timeout: got %v, want %v", client.timeout, 10*time.Second)
	}
	if client.httpClient == nil {
		t.Error("expected non-nil httpClient")
	}
	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("unexpected httpClient.Timeout: got %v, want %v", client.httpClient.Timeout, 10*time.Second)
	}
}
