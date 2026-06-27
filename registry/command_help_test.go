package registry

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout during fn execution and returns what was printed.
func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCommandHelp_ReturnsNil(t *testing.T) {
	cfg := &PokedexContext{}
	if err := commandHelp(cfg); err != nil {
		t.Errorf("expected commandHelp to return nil, got %v", err)
	}
}

func TestCommandHelp_PrintsWelcomeHeader(t *testing.T) {
	cfg := &PokedexContext{}
	output := captureStdout(func() {
		commandHelp(cfg)
	})

	if !strings.Contains(output, "Welcome to the Pokedex!") {
		t.Errorf("expected output to contain 'Welcome to the Pokedex!', got:\n%s", output)
	}

	if !strings.Contains(output, "Available commands:") {
		t.Errorf("expected output to contain 'Available commands:', got:\n%s", output)
	}
}

func TestCommandHelp_PrintsAllCommands(t *testing.T) {
	cfg := &PokedexContext{}
	output := captureStdout(func() {
		commandHelp(cfg)
	})

	commands := GetCLICommands()
	for name, cmd := range commands {
		if !strings.Contains(output, name) {
			t.Errorf("expected output to contain command name %q, got:\n%s", name, output)
		}
		if !strings.Contains(output, cmd.description) {
			t.Errorf("expected output to contain description %q for command %q, got:\n%s", cmd.description, name, output)
		}
	}
}

func TestCommandHelp_NilContextDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("commandHelp panicked with nil context: %v", r)
		}
	}()

	captureStdout(func() {
		commandHelp(nil)
	})
}
