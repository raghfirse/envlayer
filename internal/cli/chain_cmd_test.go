package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envlayer/internal/cli"
)

func writeChainEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeChainEnv: %v", err)
	}
	return p
}

func TestRunChain_ResolvesHighestPriorityLast(t *testing.T) {
	dir := t.TempDir()
	base := writeChainEnv(t, dir, "base.env", "APP_ENV=base\nPORT=8080\n")
	prod := writeChainEnv(t, dir, "prod.env", "APP_ENV=production\n")

	var buf bytes.Buffer
	if err := cli.RunChain([]string{base, prod}, "", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
}

func TestRunChain_SingleKeyLookup(t *testing.T) {
	dir := t.TempDir()
	base := writeChainEnv(t, dir, "base.env", "DB_HOST=localhost\nDB_PORT=5432\n")
	override := writeChainEnv(t, dir, "override.env", "DB_HOST=db.prod.internal\n")

	var buf bytes.Buffer
	if err := cli.RunChain([]string{base, override}, "DB_HOST", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "db.prod.internal" {
		t.Errorf("expected db.prod.internal, got %q", buf.String())
	}
}

func TestRunChain_MissingKeyReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeChainEnv(t, dir, "a.env", "X=1\n")

	var buf bytes.Buffer
	err := cli.RunChain([]string{f}, "MISSING", &buf)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRunChain_NoFilesReturnsError(t *testing.T) {
	var buf bytes.Buffer
	if err := cli.RunChain([]string{}, "", &buf); err == nil {
		t.Fatal("expected error for empty file list")
	}
}

func TestRunChain_MissingFileReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := cli.RunChain([]string{"/nonexistent/path.env"}, "", &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunChain_OutputIsSorted(t *testing.T) {
	dir := t.TempDir()
	f := writeChainEnv(t, dir, "vars.env", "ZEBRA=z\nAPPLE=a\nMIDDLE=m\n")

	var buf bytes.Buffer
	if err := cli.RunChain([]string{f}, "", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "APPLE=a" || lines[1] != "MIDDLE=m" || lines[2] != "ZEBRA=z" {
		t.Errorf("output not sorted: %v", lines)
	}
}
