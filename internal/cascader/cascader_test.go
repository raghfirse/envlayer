package cascader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envlayer/internal/cascader"
)

func writeEnv(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
}

func TestBuild_BaseOnly(t *testing.T) {
	dir := t.TempDir()
	writeEnv(t, dir, ".env", "APP=base\nDEBUG=false\n")

	res, err := cascader.Build(cascader.CascadeOptions{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP"] != "base" {
		t.Errorf("expected APP=base, got %q", res.Vars["APP"])
	}
	if len(res.Layers) != 1 || res.Layers[0].Name != "base" {
		t.Errorf("expected single base layer, got %v", res.Layers)
	}
}

func TestBuild_EnvironmentOverridesBase(t *testing.T) {
	dir := t.TempDir()
	writeEnv(t, dir, ".env", "APP=base\nDEBUG=false\nHOST=localhost\n")
	writeEnv(t, dir, ".env.production", "DEBUG=false\nHOST=prod.example.com\n")

	res, err := cascader.Build(cascader.CascadeOptions{
		Dir:          dir,
		Environments: []string{"production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "prod.example.com" {
		t.Errorf("expected HOST=prod.example.com, got %q", res.Vars["HOST"])
	}
	if res.Vars["APP"] != "base" {
		t.Errorf("expected APP=base inherited from base, got %q", res.Vars["APP"])
	}
}

func TestBuild_MultipleEnvironmentsCascade(t *testing.T) {
	dir := t.TempDir()
	writeEnv(t, dir, ".env", "KEY=base\n")
	writeEnv(t, dir, ".env.staging", "KEY=staging\nEXTRA=yes\n")
	writeEnv(t, dir, ".env.local", "KEY=local\n")

	res, err := cascader.Build(cascader.CascadeOptions{
		Dir:          dir,
		Environments: []string{"staging", "local"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "local" {
		t.Errorf("expected KEY=local (last wins), got %q", res.Vars["KEY"])
	}
	if res.Vars["EXTRA"] != "yes" {
		t.Errorf("expected EXTRA=yes from staging, got %q", res.Vars["EXTRA"])
	}
	if len(res.Layers) != 3 {
		t.Errorf("expected 3 layers, got %d", len(res.Layers))
	}
}

func TestBuild_MissingEnvironmentLayerSkipped(t *testing.T) {
	dir := t.TempDir()
	writeEnv(t, dir, ".env", "KEY=base\n")

	res, err := cascader.Build(cascader.CascadeOptions{
		Dir:          dir,
		Environments: []string{"nonexistent"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "base" {
		t.Errorf("expected KEY=base, got %q", res.Vars["KEY"])
	}
}

func TestBuild_NoFilesReturnsError(t *testing.T) {
	dir := t.TempDir()

	_, err := cascader.Build(cascader.CascadeOptions{
		Dir:          dir,
		Environments: []string{"missing"},
	})
	if err == nil {
		t.Fatal("expected error when no .env files exist")
	}
}

func TestBuild_EmptyDirReturnsError(t *testing.T) {
	_, err := cascader.Build(cascader.CascadeOptions{})
	if err == nil {
		t.Fatal("expected error for empty Dir")
	}
}
