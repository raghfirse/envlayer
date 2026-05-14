package envrender_test

import (
	"strings"
	"testing"

	"envlayer/internal/envrender"
)

func TestRender_TableFormat(t *testing.T) {
	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	opts := envrender.DefaultOptions()
	var sb strings.Builder
	if err := envrender.Render(&sb, vars, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV") || !strings.Contains(out, "production") {
		t.Errorf("expected APP_ENV=production in table output, got:\n%s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected header row with KEY, got:\n%s", out)
	}
}

func TestRender_CompactFormat(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := envrender.Options{Format: envrender.FormatCompact}
	var sb strings.Builder
	if err := envrender.Render(&sb, vars, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in compact output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in compact output, got:\n%s", out)
	}
}

func TestRender_SummaryFormat(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "", "C": "3"}
	opts := envrender.Options{Format: envrender.FormatSummary}
	var sb strings.Builder
	if err := envrender.Render(&sb, vars, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "Total keys : 3") {
		t.Errorf("expected total 3, got:\n%s", out)
	}
	if !strings.Contains(out, "Empty values: 1") {
		t.Errorf("expected 1 empty value, got:\n%s", out)
	}
}

func TestRender_MaskValues(t *testing.T) {
	vars := map[string]string{"SECRET": "s3cr3t"}
	opts := envrender.Options{Format: envrender.FormatCompact, MaskValues: true, MaskChar: "****"}
	var sb strings.Builder
	_ = envrender.Render(&sb, vars, opts)
	out := sb.String()
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected value to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "****") {
		t.Errorf("expected mask char in output, got:\n%s", out)
	}
}

func TestRender_TruncatesLongValues(t *testing.T) {
	long := strings.Repeat("x", 100)
	vars := map[string]string{"LONG": long}
	opts := envrender.Options{Format: envrender.FormatCompact, MaxValueLen: 10}
	var sb strings.Builder
	_ = envrender.Render(&sb, vars, opts)
	out := sb.String()
	if strings.Contains(out, long) {
		t.Errorf("expected value to be truncated")
	}
	if !strings.Contains(out, "...") {
		t.Errorf("expected ellipsis in truncated output")
	}
}

func TestRender_UnknownFormat_ReturnsError(t *testing.T) {
	vars := map[string]string{"X": "y"}
	opts := envrender.Options{Format: "xml"}
	var sb strings.Builder
	err := envrender.Render(&sb, vars, opts)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestRender_EmptyMap(t *testing.T) {
	opts := envrender.DefaultOptions()
	var sb strings.Builder
	if err := envrender.Render(&sb, map[string]string{}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
