package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- GetLocations ---

func locationAreaResponse(count int, next, previous *string, names []string) LocationAreaResponse {
	results := make([]LocationArea, len(names))
	for i, n := range names {
		results[i] = LocationArea{Name: n}
	}
	return LocationAreaResponse{
		Count:    count,
		Next:     next,
		Previous: previous,
		Results:  results,
	}
}

func strPtr(s string) *string { return &s }

func TestGetLocations_WithNilURL_UsesDefaultEndpoint(t *testing.T) {
	want := locationAreaResponse(2, strPtr("next-page"), nil, []string{"canalave-city-area", "eterna-city-area"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Point the default endpoint at our test server.
	original := Endpoints["location-area"]
	Endpoints["location-area"] = server.URL
	defer func() { Endpoints["location-area"] = original }()

	client := NewPokedexClient()
	got, err := client.GetLocations(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Count != want.Count {
		t.Errorf("Count: got %d, want %d", got.Count, want.Count)
	}
	if len(got.Results) != len(want.Results) {
		t.Fatalf("Results length: got %d, want %d", len(got.Results), len(want.Results))
	}
	for i, r := range got.Results {
		if r.Name != want.Results[i].Name {
			t.Errorf("Results[%d].Name: got %q, want %q", i, r.Name, want.Results[i].Name)
		}
	}
}

func TestGetLocations_WithExplicitURL_OverridesDefault(t *testing.T) {
	want := locationAreaResponse(1, nil, strPtr("prev-page"), []string{"pallet-town-area"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewPokedexClient()
	url := server.URL
	got, err := client.GetLocations(&url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Count != want.Count {
		t.Errorf("Count: got %d, want %d", got.Count, want.Count)
	}
	if len(got.Results) != 1 || got.Results[0].Name != "pallet-town-area" {
		t.Errorf("unexpected results: %+v", got.Results)
	}
	if got.Previous == nil || *got.Previous != "prev-page" {
		t.Errorf("Previous: got %v, want %q", got.Previous, "prev-page")
	}
}

func TestGetLocations_ServerError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewPokedexClient()
	url := server.URL
	// A non-200 response with a non-JSON body should cause json.Unmarshal to fail.
	_, err := client.GetLocations(&url)
	if err == nil {
		t.Error("expected error for malformed JSON from server error response, got nil")
	}
}

func TestGetLocations_InvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not valid json{{"))
	}))
	defer server.Close()

	client := NewPokedexClient()
	url := server.URL
	_, err := client.GetLocations(&url)
	if err == nil {
		t.Error("expected JSON parse error, got nil")
	}
}

func TestGetLocations_UnreachableServer_ReturnsError(t *testing.T) {
	// Use a URL that will immediately refuse connections.
	client := NewPokedexClient()
	bad := "http://127.0.0.1:1" // port 1 is never open
	_, err := client.GetLocations(&bad)
	if err == nil {
		t.Error("expected connection error for unreachable server, got nil")
	}
}

func TestGetLocations_PaginationFields(t *testing.T) {
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"
	prevURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	want := locationAreaResponse(100, &nextURL, &prevURL, []string{"area-a", "area-b"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewPokedexClient()
	url := server.URL
	got, err := client.GetLocations(&url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Next == nil || *got.Next != nextURL {
		t.Errorf("Next: got %v, want %q", got.Next, nextURL)
	}
	if got.Previous == nil || *got.Previous != prevURL {
		t.Errorf("Previous: got %v, want %q", got.Previous, prevURL)
	}
	if got.Count != 100 {
		t.Errorf("Count: got %d, want 100", got.Count)
	}
}
