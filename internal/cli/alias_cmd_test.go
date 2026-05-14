package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlayer/internal/cli"
)

func writeAliasEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeAliasEnv: %v", err)
	}
	return p
}

func TestRunAlias_RenamesKey(t *testing.T) {
	dir := t.TempDir()
	f := writeAliasEnv(t, dir, ".env", "DB_HOST=localhost\n")
	var buf bytes.Buffer
	err := cli.RunAlias(cli.AliasArgs{
		Files:   []string{f},
		Aliases: []string{"DB_HOST=DATABASE_HOST"},
		Out:     &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", result["DATABASE_HOST"])
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("original key DB_HOST should not appear in output")
	}
}

func TestRunAlias_KeepOriginal(t *testing.T) {
	dir := t.TempDir()
	f := writeAliasEnv(t, dir, ".env", "PORT=9000\n")
	var buf bytes.Buffer
	err := cli.RunAlias(cli.AliasArgs{
		Files:        []string{f},
		Aliases:      []string{"PORT=APP_PORT"},
		KeepOriginal: true,
		Out:          &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	json.Unmarshal(buf.Bytes(), &result)
	if result["PORT"] != "9000" || result["APP_PORT"] != "9000" {
		t.Errorf("both original and alias should be present: %v", result)
	}
}

func TestRunAlias_NoFiles_ReturnsError(t *testing.T) {
	err := cli.RunAlias(cli.AliasArgs{
		Aliases: []string{"A=B"},
	})
	if err == nil {
		t.Fatal("expected error for no files")
	}
}

func TestRunAlias_NoAliases_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeAliasEnv(t, dir, ".env", "X=1\n")
	err := cli.RunAlias(cli.AliasArgs{Files: []string{f}})
	if err == nil {
		t.Fatal("expected error for no aliases")
	}
}

func TestRunAlias_InvalidMapping_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeAliasEnv(t, dir, ".env", "X=1\n")
	err := cli.RunAlias(cli.AliasArgs{
		Files:   []string{f},
		Aliases: []string{"BADFORMAT"},
	})
	if err == nil {
		t.Fatal("expected error for malformed alias")
	}
}

func TestRunAlias_FailOnMissing_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := writeAliasEnv(t, dir, ".env", "EXISTING=yes\n")
	var buf bytes.Buffer
	err := cli.RunAlias(cli.AliasArgs{
		Files:         []string{f},
		Aliases:       []string{"MISSING=ALIAS"},
		FailOnMissing: true,
		Out:           &buf,
	})
	if err == nil {
		t.Fatal("expected error when source key is missing")
	}
}
