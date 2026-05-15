package envpinner_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envpinner"
)

func TestPin_PinnedKeyNotOverridden(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	incoming := map[string]string{"HOST": "remotehost", "PORT": "9999"}

	res, err := envpinner.Pin(base, incoming, []string{"HOST"}, envpinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Vars["HOST"])
	}
	if res.Vars["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %q", res.Vars["PORT"])
	}
}

func TestPin_ViolationRecorded(t *testing.T) {
	base := map[string]string{"SECRET": "original"}
	incoming := map[string]string{"SECRET": "changed"}

	res, err := envpinner.Pin(base, incoming, []string{"SECRET"}, envpinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Violations) != 1 || res.Violations[0] != "SECRET" {
		t.Errorf("expected violation for SECRET, got %v", res.Violations)
	}
}

func TestPin_StrictMode_ReturnsError(t *testing.T) {
	base := map[string]string{"API_KEY": "abc"}
	incoming := map[string]string{"API_KEY": "xyz"}
	opts := envpinner.Options{StrictMode: true}

	_, err := envpinner.Pin(base, incoming, []string{"API_KEY"}, opts)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestPin_StrictMode_SameValueNoError(t *testing.T) {
	base := map[string]string{"API_KEY": "abc"}
	incoming := map[string]string{"API_KEY": "abc"}
	opts := envpinner.Options{StrictMode: true}

	res, err := envpinner.Pin(base, incoming, []string{"API_KEY"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["API_KEY"] != "abc" {
		t.Errorf("expected API_KEY=abc, got %q", res.Vars["API_KEY"])
	}
}

func TestPin_NewKeyInIncomingAdded(t *testing.T) {
	base := map[string]string{"A": "1"}
	incoming := map[string]string{"B": "2"}

	res, err := envpinner.Pin(base, incoming, []string{"A"}, envpinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["B"] != "2" {
		t.Errorf("expected B=2, got %q", res.Vars["B"])
	}
}

func TestPin_SilentMode_NoViolations(t *testing.T) {
	base := map[string]string{"TOKEN": "secret"}
	incoming := map[string]string{"TOKEN": "other"}
	opts := envpinner.Options{Silent: true}

	res, err := envpinner.Pin(base, incoming, []string{"TOKEN"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Violations) != 0 {
		t.Errorf("expected no violations in silent mode, got %v", res.Violations)
	}
}

func TestPinnedKeys_ReturnsPresent(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	keys := envpinner.PinnedKeys(vars, []string{"A", "D", "C"})
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "C" {
		t.Errorf("unexpected pinned keys: %v", keys)
	}
}
