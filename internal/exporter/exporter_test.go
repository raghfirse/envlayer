package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlayer/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"DB_PORT":  "5432",
	"SECRET":   "p@ss\"word",
}

func TestExport_ShellFormat(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"APP_ENV", "DB_HOST", "DB_PORT", "SECRET"} {
		if !strings.Contains(out, "export "+key+"=") {
			t.Errorf("expected 'export %s=' in output, got:\n%s", key, out)
		}
	}
}

func TestExport_ShellFormat_SortedOutput(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(sampleEnv) {
		t.Fatalf("expected %d lines, got %d", len(sampleEnv), len(lines))
	}
	// First key alphabetically should be APP_ENV.
	if !strings.HasPrefix(lines[0], "export APP_ENV=") {
		t.Errorf("expected first line to start with 'export APP_ENV=', got: %s", lines[0])
	}
}

func TestExport_DotenvFormat(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("dotenv format should not contain 'export' keyword")
	}
	if !strings.Contains(out, "DB_PORT=") {
		t.Errorf("expected 'DB_PORT=' in dotenv output")
	}
}

func TestExport_JSONFormat(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	for key := range sampleEnv {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in JSON output", key)
		}
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	_, err := exporter.Export(sampleEnv, exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestExport_EmptyMap(t *testing.T) {
	for _, fmt := range []exporter.Format{exporter.FormatExport, exporter.FormatDotenv, exporter.FormatJSON} {
		out, err := exporter.Export(map[string]string{}, fmt)
		if err != nil {
			t.Errorf("format %q: unexpected error: %v", fmt, err)
		}
		if fmt == exporter.FormatJSON && !strings.Contains(out, "{") {
			t.Errorf("format %q: expected at least an empty JSON object", fmt)
		}
	}
}
