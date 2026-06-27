package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thalesraymond/pokedexcli/api"
)

// helpers shared across map tests

func strPtr(s string) *string { return &s }

func serveLocationAreaResponse(t *testing.T, resp api.LocationAreaResponse) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
}

func overrideEndpoint(t *testing.T, url string) {
	t.Helper()
	original := api.Endpoints["location-area"]
	api.Endpoints["location-area"] = url
	t.Cleanup(func() { api.Endpoints["location-area"] = original })
}

// --- commandMap tests ---

func TestCommandMap_PrintsLocationNames(t *testing.T) {
	want := api.LocationAreaResponse{
		Count: 2,
		Next:  strPtr("http://example.com/next"),
		Results: []api.LocationAreaDto{
			{Name: "canalave-city-area"},
			{Name: "eterna-city-area"},
		},
	}

	server := serveLocationAreaResponse(t, want)
	defer server.Close()
	overrideEndpoint(t, server.URL)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandMap_UpdatesNextURL(t *testing.T) {
	nextURL := "http://example.com/next-page"
	want := api.LocationAreaResponse{
		Count: 1,
		Next:  &nextURL,
		Results: []api.LocationAreaDto{{Name: "pallet-town-area"}},
	}

	server := serveLocationAreaResponse(t, want)
	defer server.Close()
	overrideEndpoint(t, server.URL)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LocationAreasNextURL == nil {
		t.Fatal("expected LocationAreasNextURL to be set, got nil")
	}
	if *cfg.LocationAreasNextURL != nextURL {
		t.Errorf("LocationAreasNextURL: got %q, want %q", *cfg.LocationAreasNextURL, nextURL)
	}
}

func TestCommandMap_UpdatesPreviousURL(t *testing.T) {
	prevURL := "http://example.com/prev-page"
	want := api.LocationAreaResponse{
		Count:    1,
		Previous: &prevURL,
		Results:  []api.LocationAreaDto{{Name: "viridian-city-area"}},
	}

	server := serveLocationAreaResponse(t, want)
	defer server.Close()
	overrideEndpoint(t, server.URL)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LocationAreasPreviousURL == nil {
		t.Fatal("expected LocationAreasPreviousURL to be set, got nil")
	}
	if *cfg.LocationAreasPreviousURL != prevURL {
		t.Errorf("LocationAreasPreviousURL: got %q, want %q", *cfg.LocationAreasPreviousURL, prevURL)
	}
}

func TestCommandMap_NilNextAndPrevious_DoesNotUpdateContext(t *testing.T) {
	want := api.LocationAreaResponse{
		Count:   1,
		Results: []api.LocationAreaDto{{Name: "cerulean-city-area"}},
	}

	server := serveLocationAreaResponse(t, want)
	defer server.Close()
	overrideEndpoint(t, server.URL)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LocationAreasNextURL != nil {
		t.Errorf("expected LocationAreasNextURL to remain nil, got %q", *cfg.LocationAreasNextURL)
	}
	if cfg.LocationAreasPreviousURL != nil {
		t.Errorf("expected LocationAreasPreviousURL to remain nil, got %q", *cfg.LocationAreasPreviousURL)
	}
}

func TestCommandMap_UsesExplicitNextURLFromContext(t *testing.T) {
	// Simulate having already paged forward: context holds a next URL.
	want := api.LocationAreaResponse{
		Count:   1,
		Results: []api.LocationAreaDto{{Name: "lavender-town-area"}},
	}

	server := serveLocationAreaResponse(t, want)
	defer server.Close()

	page2URL := server.URL
	// Do NOT override the default endpoint — cfg carries the explicit URL.
	cfg := &PokedexContext{LocationAreasNextURL: &page2URL}
	if err := commandMap(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandMap_NetworkError_ReturnsError(t *testing.T) {
	bad := "http://127.0.0.1:1" // port 1 is never open
	overrideEndpoint(t, bad)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}

func TestCommandMap_InvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not valid json{{"))
	}))
	defer server.Close()
	overrideEndpoint(t, server.URL)

	cfg := &PokedexContext{}
	if err := commandMap(cfg); err == nil {
		t.Error("expected JSON parse error, got nil")
	}
}
