package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func makeEnvDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}
	return dir
}

func TestRun_BasicDotenv(t *testing.T) {
	dir := makeEnvDir(t, map[string]string{
		".env": "APP=hello\nDEBUG=false\n",
	})
	code := Run([]string{"-dir", dir, "-format", "dotenv"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_WithEnvironmentOverride(t *testing.T) {
	dir := makeEnvDir(t, map[string]string{
		".env":            "APP=base\nPORT=3000\n",
		".env.staging":    "PORT=4000\n",
	})
	code := Run([]string{"-dir", dir, "-env", "staging", "-format", "dotenv"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_MissingRequiredKey(t *testing.T) {
	dir := makeEnvDir(t, map[string]string{
		".env": "APP=hello\n",
	})
	code := Run([]string{"-dir", dir, "-require", "APP,SECRET_KEY"})
	if code == 0 {
		t.Error("expected non-zero exit code when required key is missing")
	}
}

func TestRun_AllRequiredPresent(t *testing.T) {
	dir := makeEnvDir(t, map[string]string{
		".env": "APP=hello\nSECRET_KEY=abc123\n",
	})
	code := Run([]string{"-dir", dir, "-require", "APP,SECRET_KEY", "-format", "dotenv"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_NoEnvFilesReturnsError(t *testing.T) {
	dir := t.TempDir()
	code := Run([]string{"-dir", dir})
	if code == 0 {
		t.Error("expected non-zero exit code when no .env files found")
	}
}

func TestRun_ExportFlag(t *testing.T) {
	dir := makeEnvDir(t, map[string]string{
		".env": "EXPORTED=yes\n",
	})
	code := Run([]string{"-dir", dir, "-export", "-format", "shell"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestSplitCSV(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"A,B,C", []string{"A", "B", "C"}},
		{"SINGLE", []string{"SINGLE"}},
		{"", []string{}},
		{"X,,Y", []string{"X", "Y"}},
	}
	for _, tc := range cases {
		got := splitCSV(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("splitCSV(%q): got %v, want %v", tc.input, got, tc.expected)
			continue
		}
		for i := range got {
			if got[i] != tc.expected[i] {
				t.Errorf("splitCSV(%q)[%d]: got %q, want %q", tc.input, i, got[i], tc.expected[i])
			}
		}
	}
}
