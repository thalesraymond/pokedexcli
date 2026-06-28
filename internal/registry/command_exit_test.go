package registry

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCommandExit_PrintsGoodbyeAndExitsZero verifies that commandExit prints
// the expected farewell message and terminates the process with exit code 0.
//
// Because commandExit calls os.Exit directly, we use the subprocess pattern:
// re-invoke the test binary with a special env flag so that only the target
// code runs inside the subprocess, then assert its output and exit code from
// the parent process.
func TestCommandExit_PrintsGoodbyeAndExitsZero(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestCommandExitSubprocess", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_COMMAND_EXIT_SUBPROCESS=1")

	out, err := cmd.CombinedOutput()

	// The subprocess must exit with code 0 (os.Exit(0) inside commandExit).
	if err != nil {
		t.Fatalf("expected exit code 0, got error: %v\nsubprocess output:\n%s", err, out)
	}

	output := string(out)
	want := "Closing the Pokedex... Goodbye!"
	if !strings.Contains(output, want) {
		t.Errorf("expected subprocess output to contain %q\ngot:\n%s", want, output)
	}
}

// TestCommandExitSubprocess is the subprocess entry point used by
// TestCommandExit_PrintsGoodbyeAndExitsZero. It is skipped unless the process
// is running as the subprocess (TEST_COMMAND_EXIT_SUBPROCESS=1).
func TestCommandExitSubprocess(t *testing.T) {
	if os.Getenv("TEST_COMMAND_EXIT_SUBPROCESS") != "1" {
		t.Skip("subprocess helper – only runs inside the test subprocess")
	}

	ctx := &PokedexContext{}
	// commandExit calls os.Exit(0), which terminates this subprocess process.
	_ = commandExit(ctx)
}
