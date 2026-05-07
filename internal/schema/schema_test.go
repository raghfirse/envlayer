package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlayer/internal/schema"
)

func writeSchema(t *testing.T, s schema.Schema) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadSchema_Valid(t *testing.T) {
	s := schema.Schema{
		"PORT": {Type: schema.TypeInt, Required: true, Description: "HTTP port"},
	}
	path := writeSchema(t, s)
	loaded, err := schema.LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := loaded["PORT"]; !ok {
		t.Error("expected PORT in loaded schema")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := schema.LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	s := schema.Schema{
		"PORT": {Type: schema.TypeInt, Required: true},
	}
	violations, err := schema.Validate(s, map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatal(err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	s := schema.Schema{
		"SECRET": {Type: schema.TypeString, Required: true},
	}
	violations, _ := schema.Validate(s, map[string]string{})
	if len(violations) != 1 || violations[0].Key != "SECRET" {
		t.Errorf("expected violation for SECRET, got %v", violations)
	}
}

func TestValidate_WrongType_Int(t *testing.T) {
	s := schema.Schema{
		"PORT": {Type: schema.TypeInt, Required: true},
	}
	violations, _ := schema.Validate(s, map[string]string{"PORT": "not-a-number"})
	if len(violations) != 1 {
		t.Errorf("expected 1 type violation, got %v", violations)
	}
}

func TestValidate_WrongType_Bool(t *testing.T) {
	s := schema.Schema{
		"DEBUG": {Type: schema.TypeBool, Required: true},
	}
	violations, _ := schema.Validate(s, map[string]string{"DEBUG": "yes-please"})
	if len(violations) != 1 {
		t.Errorf("expected 1 type violation, got %v", violations)
	}
}

func TestValidate_DefaultSkipsMissing(t *testing.T) {
	s := schema.Schema{
		"TIMEOUT": {Type: schema.TypeInt, Required: true, Default: "30"},
	}
	violations, _ := schema.Validate(s, map[string]string{})
	if len(violations) != 0 {
		t.Errorf("expected no violations when default is set, got %v", violations)
	}
}

func TestValidate_FloatType(t *testing.T) {
	s := schema.Schema{
		"RATE": {Type: schema.TypeFloat, Required: true},
	}
	violations, _ := schema.Validate(s, map[string]string{"RATE": "3.14"})
	if len(violations) != 0 {
		t.Errorf("expected no violations for valid float, got %v", violations)
	}
}
