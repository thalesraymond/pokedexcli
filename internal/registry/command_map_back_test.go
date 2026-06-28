package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thalesraymond/pokedexcli/internal/api"
)

// buildLocationResponse is a convenience builder for LocationAreaResponse test fixtures.
func buildLocationResponse(next, previous *string, names []string) api.LocationAreaResponse {
	results := make([]api.LocationAreaDto, len(names))
	for i, n := range names {
		results[i] = api.LocationAreaDto{Name: n}
	}
	return api.LocationAreaResponse{
		Count:    len(names),
		Next:     next,
		Previous: previous,
		Results:  results,
	}
}

// newTestServer spins up an httptest.Server that always responds with resp encoded as JSON.
func newTestServer(t *testing.T, resp api.LocationAreaResponse) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("test server encode error: %v", err)
		}
	}))
}

// TestCommandMapBack_NilPreviousURL verifies that when LocationAreasPreviousURL
// is nil the function prints the "first page" message and returns nil (no HTTP call).
func TestCommandMapBack_NilPreviousURL(t *testing.T) {
	cfg := &PokedexContext{} // both URL fields are nil

	output := captureStdout(func() {
		if err := commandMapBack(cfg); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	if !strings.Contains(output, "you're on the first page") {
		t.Errorf("expected \"you're on the first page\" in output, got:\n%s", output)
	}
}

// TestCommandMapBack_PrintsLocationNames verifies that each location name returned
// by the API is printed to stdout.
func TestCommandMapBack_PrintsLocationNames(t *testing.T) {
	wantNames := []string{"canalave-city-area", "eterna-city-area", "pastoria-city-area"}
	resp := buildLocationResponse(nil, nil, wantNames)

	server := newTestServer(t, resp)
	defer server.Close()

	prevURL := server.URL
	cfg := &PokedexContext{LocationAreasPreviousURL: &prevURL}

	output := captureStdout(func() {
		if err := commandMapBack(cfg); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	for _, name := range wantNames {
		if !strings.Contains(output, name) {
			t.Errorf("expected output to contain %q, got:\n%s", name, output)
		}
	}
}

// TestCommandMapBack_UpdatesNextURL verifies that cfg.LocationAreasNextURL is
// updated when the API response contains a non-nil Next field.
func TestCommandMapBack_UpdatesNextURL(t *testing.T) {
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	resp := buildLocationResponse(&nextURL, nil, []string{"area-a"})

	server := newTestServer(t, resp)
	defer server.Close()

	prevURL := server.URL
	cfg := &PokedexContext{LocationAreasPreviousURL: &prevURL}

	captureStdout(func() {
		_ = commandMapBack(cfg)
	})

	if cfg.LocationAreasNextURL == nil {
		t.Fatal("expected LocationAreasNextURL to be set, got nil")
	}
	if *cfg.LocationAreasNextURL != nextURL {
		t.Errorf("LocationAreasNextURL: got %q, want %q", *cfg.LocationAreasNextURL, nextURL)
	}
}

// TestCommandMapBack_UpdatesPreviousURL verifies that cfg.LocationAreasPreviousURL
// is updated when the API response contains a non-nil Previous field.
func TestCommandMapBack_UpdatesPreviousURL(t *testing.T) {
	newPrevURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	resp := buildLocationResponse(nil, &newPrevURL, []string{"area-b"})

	server := newTestServer(t, resp)
	defer server.Close()

	prevURL := server.URL
	cfg := &PokedexContext{LocationAreasPreviousURL: &prevURL}

	captureStdout(func() {
		_ = commandMapBack(cfg)
	})

	if cfg.LocationAreasPreviousURL == nil {
		t.Fatal("expected LocationAreasPreviousURL to be set, got nil")
	}
	if *cfg.LocationAreasPreviousURL != newPrevURL {
		t.Errorf("LocationAreasPreviousURL: got %q, want %q", *cfg.LocationAreasPreviousURL, newPrevURL)
	}
}

// TestCommandMapBack_DoesNotUpdateNextURLWhenNil verifies that
// cfg.LocationAreasNextURL is left unchanged when the API returns Next == nil.
func TestCommandMapBack_DoesNotUpdateNextURLWhenNil(t *testing.T) {
	originalNext := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
	resp := buildLocationResponse(nil, nil, []string{"area-c"})

	server := newTestServer(t, resp)
	defer server.Close()

	prevURL := server.URL
	cfg := &PokedexContext{
		LocationAreasNextURL:     &originalNext,
		LocationAreasPreviousURL: &prevURL,
	}

	captureStdout(func() {
		_ = commandMapBack(cfg)
	})

	if cfg.LocationAreasNextURL == nil || *cfg.LocationAreasNextURL != originalNext {
		t.Errorf("LocationAreasNextURL should be unchanged: got %v, want %q",
			cfg.LocationAreasNextURL, originalNext)
	}
}

// TestCommandMapBack_ReturnsErrorOnBadURL verifies that a network failure
// (unreachable URL) causes commandMapBack to return a non-nil error.
func TestCommandMapBack_ReturnsErrorOnBadURL(t *testing.T) {
	bad := "http://127.0.0.1:1" // port 1 is never open
	cfg := &PokedexContext{LocationAreasPreviousURL: &bad}

	var err error
	captureStdout(func() {
		err = commandMapBack(cfg)
	})

	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}

// TestCommandMapBack_ReturnsNilOnSuccess verifies the happy-path return value.
func TestCommandMapBack_ReturnsNilOnSuccess(t *testing.T) {
	resp := buildLocationResponse(nil, nil, []string{"viridian-city-area"})

	server := newTestServer(t, resp)
	defer server.Close()

	prevURL := server.URL
	cfg := &PokedexContext{LocationAreasPreviousURL: &prevURL}

	captureStdout(func() {
		if err := commandMapBack(cfg); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}
