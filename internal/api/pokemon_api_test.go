package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- helpers ---

func buildPokemon(id int, name string, baseExp, height, weight int, types []string) Pokemon {
	t := make([]Types, len(types))
	for i, typeName := range types {
		t[i] = Types{Slot: i + 1, Type: Type{Name: typeName}}
	}
	return Pokemon{
		ID:             id,
		Name:           name,
		BaseExperience: baseExp,
		Height:         height,
		Weight:         weight,
		Types:          t,
	}
}

func servePokemon(t *testing.T, p Pokemon) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(p); err != nil {
			t.Errorf("failed to encode pokemon response: %v", err)
		}
	}))
}

// redirectEndpoint temporarily replaces the "pokemon" endpoint and returns a
// restore function to defer.
func redirectPokemonEndpoint(url string) func() {
	original := Endpoints["pokemon"]
	Endpoints["pokemon"] = url
	return func() { Endpoints["pokemon"] = original }
}

// --- GetPokemon ---

func TestGetPokemon_HappyPath_ReturnsCorrectFields(t *testing.T) {
	want := buildPokemon(6, "charizard", 240, 17, 905, []string{"fire", "flying"})

	server := servePokemon(t, want)
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	got, err := client.GetPokemon("charizard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != want.ID {
		t.Errorf("ID: got %d, want %d", got.ID, want.ID)
	}
	if got.Name != want.Name {
		t.Errorf("Name: got %q, want %q", got.Name, want.Name)
	}
	if got.BaseExperience != want.BaseExperience {
		t.Errorf("BaseExperience: got %d, want %d", got.BaseExperience, want.BaseExperience)
	}
	if got.Height != want.Height {
		t.Errorf("Height: got %d, want %d", got.Height, want.Height)
	}
	if got.Weight != want.Weight {
		t.Errorf("Weight: got %d, want %d", got.Weight, want.Weight)
	}
	if len(got.Types) != len(want.Types) {
		t.Fatalf("Types length: got %d, want %d", len(got.Types), len(want.Types))
	}
	for i, tp := range got.Types {
		if tp.Type.Name != want.Types[i].Type.Name {
			t.Errorf("Types[%d].Type.Name: got %q, want %q", i, tp.Type.Name, want.Types[i].Type.Name)
		}
	}
}

func TestGetPokemon_InvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{{not valid json"))
	}))
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	_, err := client.GetPokemon("pikachu")
	if err == nil {
		t.Error("expected JSON parse error, got nil")
	}
}

func TestGetPokemon_UnreachableServer_ReturnsError(t *testing.T) {
	defer redirectPokemonEndpoint("http://127.0.0.1:1")()

	client := NewPokedexClient()
	_, err := client.GetPokemon("bulbasaur")
	if err == nil {
		t.Error("expected connection error for unreachable server, got nil")
	}
}

func TestGetPokemon_ServerError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	_, err := client.GetPokemon("mewtwo")
	if err == nil {
		t.Error("expected error for non-JSON server error response, got nil")
	}
}

func TestGetPokemon_CachesResponse(t *testing.T) {
	callCount := 0
	want := buildPokemon(25, "pikachu", 112, 4, 60, []string{"electric"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()

	// First call — should hit the server.
	got1, err := client.GetPokemon("pikachu")
	if err != nil {
		t.Fatalf("first call unexpected error: %v", err)
	}

	// Second call — should be served from cache.
	got2, err := client.GetPokemon("pikachu")
	if err != nil {
		t.Fatalf("second call unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected exactly 1 HTTP request, got %d", callCount)
	}
	if got1.Name != got2.Name {
		t.Errorf("cached result differs: got1=%q, got2=%q", got1.Name, got2.Name)
	}
}

func TestGetPokemon_DifferentPokemonNames_HitServerIndependently(t *testing.T) {
	pikachu := buildPokemon(25, "pikachu", 112, 4, 60, []string{"electric"})
	eevee := buildPokemon(133, "eevee", 65, 3, 65, []string{"normal"})

	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		var p Pokemon
		switch r.URL.Path {
		case "/pikachu":
			p = pikachu
		case "/eevee":
			p = eevee
		default:
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(p); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()

	gotPika, err := client.GetPokemon("pikachu")
	if err != nil {
		t.Fatalf("pikachu: unexpected error: %v", err)
	}
	gotEevee, err := client.GetPokemon("eevee")
	if err != nil {
		t.Fatalf("eevee: unexpected error: %v", err)
	}

	if gotPika.Name != "pikachu" {
		t.Errorf("expected pikachu, got %q", gotPika.Name)
	}
	if gotEevee.Name != "eevee" {
		t.Errorf("expected eevee, got %q", gotEevee.Name)
	}
	if callCount != 2 {
		t.Errorf("expected 2 HTTP requests (one per pokemon), got %d", callCount)
	}
}

func TestGetPokemon_Stats_DeserializedCorrectly(t *testing.T) {
	want := Pokemon{
		ID:   94,
		Name: "gengar",
		Stats: []Stats{
			{BaseStat: 60, Effort: 0, Stat: Stat{Name: "hp"}},
			{BaseStat: 65, Effort: 0, Stat: Stat{Name: "attack"}},
			{BaseStat: 130, Effort: 3, Stat: Stat{Name: "special-attack"}},
		},
	}

	server := servePokemon(t, want)
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	got, err := client.GetPokemon("gengar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Stats) != len(want.Stats) {
		t.Fatalf("Stats length: got %d, want %d", len(got.Stats), len(want.Stats))
	}
	for i, s := range got.Stats {
		wantStat := want.Stats[i]
		if s.BaseStat != wantStat.BaseStat {
			t.Errorf("Stats[%d].BaseStat: got %d, want %d", i, s.BaseStat, wantStat.BaseStat)
		}
		if s.Effort != wantStat.Effort {
			t.Errorf("Stats[%d].Effort: got %d, want %d", i, s.Effort, wantStat.Effort)
		}
		if s.Stat.Name != wantStat.Stat.Name {
			t.Errorf("Stats[%d].Stat.Name: got %q, want %q", i, s.Stat.Name, wantStat.Stat.Name)
		}
	}
}

func TestGetPokemon_Abilities_DeserializedCorrectly(t *testing.T) {
	want := Pokemon{
		ID:   6,
		Name: "charizard",
		Abilities: []Abilities{
			{IsHidden: false, Slot: 1, Ability: Ability{Name: "blaze", URL: "https://pokeapi.co/api/v2/ability/66/"}},
			{IsHidden: true, Slot: 3, Ability: Ability{Name: "solar-power", URL: "https://pokeapi.co/api/v2/ability/94/"}},
		},
	}

	server := servePokemon(t, want)
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	got, err := client.GetPokemon("charizard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Abilities) != len(want.Abilities) {
		t.Fatalf("Abilities length: got %d, want %d", len(got.Abilities), len(want.Abilities))
	}
	for i, a := range got.Abilities {
		wa := want.Abilities[i]
		if a.Ability.Name != wa.Ability.Name {
			t.Errorf("Abilities[%d].Ability.Name: got %q, want %q", i, a.Ability.Name, wa.Ability.Name)
		}
		if a.IsHidden != wa.IsHidden {
			t.Errorf("Abilities[%d].IsHidden: got %v, want %v", i, a.IsHidden, wa.IsHidden)
		}
		if a.Slot != wa.Slot {
			t.Errorf("Abilities[%d].Slot: got %d, want %d", i, a.Slot, wa.Slot)
		}
	}
}

func TestGetPokemon_EmptyName_StillCallsEndpoint(t *testing.T) {
	want := buildPokemon(0, "", 0, 0, 0, nil)

	server := servePokemon(t, want)
	defer server.Close()

	defer redirectPokemonEndpoint(server.URL)()

	client := NewPokedexClient()
	// URL becomes "<endpoint>/", the server still responds — verify no error.
	_, err := client.GetPokemon("")
	if err != nil {
		t.Fatalf("unexpected error for empty name: %v", err)
	}
}
