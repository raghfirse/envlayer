package redactor_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/redactor"
)

func TestRedact_SensitiveValueReplaced(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
	}
	r := redactor.New(vars, "")
	out := r.Redact("connecting with supersecret to myapp")
	if got, want := out, "connecting with **** to myapp"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_NonSensitiveValueNotReplaced(t *testing.T) {
	vars := map[string]string{
		"APP_ENV": "production",
	}
	r := redactor.New(vars, "")
	out := r.Redact("running in production mode")
	if got, want := out, "running in production mode"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_CustomMaskChar(t *testing.T) {
	vars := map[string]string{
		"API_SECRET": "topsecret",
	}
	r := redactor.New(vars, "[REDACTED]")
	out := r.Redact("token=topsecret")
	if got, want := out, "token=[REDACTED]"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_EmptyVarsNoChange(t *testing.T) {
	r := redactor.New(map[string]string{}, "")
	out := r.Redact("nothing to hide here")
	if got, want := out, "nothing to hide here"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedactLines_MultiLine(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "hunter2",
	}
	r := redactor.New(vars, "***")
	input := "line1\npassword=hunter2\nline3"
	out := r.RedactLines(input)
	want := "line1\npassword=***\nline3"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestRedact_EmptyValueSkipped(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "",
	}
	r := redactor.New(vars, "")
	// empty string replacement would corrupt all output; ensure it is skipped
	out := r.Redact("some output text")
	if got, want := out, "some output text"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
