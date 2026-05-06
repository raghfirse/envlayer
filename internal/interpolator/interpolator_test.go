package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_BracedSyntax(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello",
		"MESSAGE":  "${GREETING}, world!",
	}
	got, err := Interpolate(env, nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["MESSAGE"] != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got %q", got["MESSAGE"])
	}
}

func TestInterpolate_UnbracedSyntax(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://$HOST:8080",
	}
	got, err := Interpolate(env, nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["URL"] != "http://localhost:8080" {
		t.Errorf("expected 'http://localhost:8080', got %q", got["URL"])
	}
}

func TestInterpolate_FallsBackToBase(t *testing.T) {
	base := map[string]string{"REGION": "us-east-1"}
	env := map[string]string{"BUCKET": "my-bucket-${REGION}"}
	got, err := Interpolate(env, base, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["BUCKET"] != "my-bucket-us-east-1" {
		t.Errorf("got %q", got["BUCKET"])
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "fromOS")
	_ = os.Getenv("OS_VAR") // ensure set
	env := map[string]string{"VAL": "${OS_VAR}_suffix"}
	got, err := Interpolate(env, nil, Options{FallbackToOS: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["VAL"] != "fromOS_suffix" {
		t.Errorf("got %q", got["VAL"])
	}
}

func TestInterpolate_MissingVar_EmptyString(t *testing.T) {
	env := map[string]string{"A": "${MISSING}"}
	got, err := Interpolate(env, nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "" {
		t.Errorf("expected empty string, got %q", got["A"])
	}
}

func TestInterpolate_ErrorOnMissing(t *testing.T) {
	env := map[string]string{"A": "${UNDEFINED}"}
	_, err := Interpolate(env, nil, Options{ErrorOnMissing: true})
	if err == nil {
		t.Fatal("expected error for undefined variable")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"X": "${Y}", "Y": "hello"}
	original := env["X"]
	_, _ = Interpolate(env, nil, Options{})
	if env["X"] != original {
		t.Error("input map was mutated")
	}
}
