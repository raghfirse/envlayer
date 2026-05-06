package templater_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envlayer/internal/templater"
)

func TestRender_BasicSubstitution(t *testing.T) {
	vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	out, err := templater.Render("{{.APP_HOST}}:{{.APP_PORT}}", vars, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "localhost:8080" {
		t.Errorf("expected 'localhost:8080', got %q", out)
	}
}

func TestRender_MultiLineTemplate(t *testing.T) {
	tmpl := "HOST={{.DB_HOST}}\nPORT={{.DB_PORT}}\nNAME={{.DB_NAME}}"
	vars := map[string]string{"DB_HOST": "db.local", "DB_PORT": "5432", "DB_NAME": "mydb"}
	out, err := templater.Render(tmpl, vars, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HOST=db.local") || !strings.Contains(out, "PORT=5432") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	vars := map[string]string{"PRESENT": "yes"}
	_, err := templater.Render("{{.MISSING}}", vars, templater.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if !strings.Contains(err.Error(), "render error") {
		t.Errorf("expected 'render error' in message, got: %v", err)
	}
}

func TestRender_MissingKey_ZeroOption(t *testing.T) {
	vars := map[string]string{}
	opts := templater.Options{MissingKey: "zero"}
	out, err := templater.Render("value={{.UNDEFINED}}", vars, opts)
	if err != nil {
		t.Fatalf("unexpected error with zero option: %v", err)
	}
	if out != "value=" {
		t.Errorf("expected empty substitution, got %q", out)
	}
}

func TestRender_InvalidTemplate_ReturnsParseError(t *testing.T) {
	_, err := templater.Render("{{.UNCLOSED", map[string]string{}, templater.DefaultOptions())
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "parse error") {
		t.Errorf("expected 'parse error' in message, got: %v", err)
	}
}

func TestRenderFile_CallsWriteWithResult(t *testing.T) {
	vars := map[string]string{"ENV": "production"}
	var captured string
	err := templater.RenderFile("env={{.ENV}}", vars, templater.DefaultOptions(), func(s string) error {
		captured = s
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured != "env=production" {
		t.Errorf("expected 'env=production', got %q", captured)
	}
}

func TestRender_EmptyTemplate(t *testing.T) {
	out, err := templater.Render("", map[string]string{"X": "y"}, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}
