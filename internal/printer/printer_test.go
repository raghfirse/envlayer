package printer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlayer/internal/printer"
)

func TestPrint_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	p := printer.New(printer.Options{Format: "shell", Output: &buf})

	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := p.Print(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "export FOO=bar") {
		t.Errorf("expected export FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "export BAZ=qux") {
		t.Errorf("expected export BAZ=qux in output, got: %s", out)
	}
}

func TestPrint_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	p := printer.New(printer.Options{Format: "dotenv", Output: &buf})

	env := map[string]string{"KEY": "value"}
	if err := p.Print(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "KEY=value") {
		t.Errorf("expected KEY=value in dotenv output, got: %s", buf.String())
	}
}

func TestPrint_WithPrefix_FiltersKeys(t *testing.T) {
	var buf bytes.Buffer
	p := printer.New(printer.Options{Format: "shell", Output: &buf, Prefix: "APP_"})

	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	if err := p.Print(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "APP_HOST") {
		t.Errorf("expected APP_HOST in output")
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Errorf("expected APP_PORT in output")
	}
	if strings.Contains(out, "DB_HOST") {
		t.Errorf("did not expect DB_HOST in prefix-filtered output")
	}
}

func TestPrint_DefaultsToStdoutFormat(t *testing.T) {
	var buf bytes.Buffer
	p := printer.New(printer.Options{Output: &buf})

	env := map[string]string{"X": "1"}
	if err := p.Print(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "export X=1") {
		t.Errorf("expected default shell format, got: %s", buf.String())
	}
}

func TestPrint_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	p := printer.New(printer.Options{Format: "xml", Output: &buf})

	env := map[string]string{"A": "b"}
	if err := p.Print(env); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
