// Package schema provides validation of environment variable definitions
// against a declared schema, supporting types, defaults, and descriptions.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// FieldType represents the expected type of an environment variable.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeFloat   FieldType = "float"
)

// Field describes a single environment variable in the schema.
type Field struct {
	Type        FieldType `json:"type"`
	Required    bool      `json:"required"`
	Default     string    `json:"default"`
	Description string    `json:"description"`
}

// Schema maps variable names to their field definitions.
type Schema map[string]Field

// Violation describes a schema validation failure.
type Violation struct {
	Key     string
	Message string
}

// LoadSchema reads a JSON schema file from the given path.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read file: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse JSON: %w", err)
	}
	return s, nil
}

// Validate checks the provided vars map against the schema.
// It returns a list of violations (may be empty) and any hard error.
func Validate(s Schema, vars map[string]string) ([]Violation, error) {
	var violations []Violation

	for key, field := range s {
		val, exists := vars[key]
		if !exists || val == "" {
			if field.Default != "" {
				continue
			}
			if field.Required {
				violations = append(violations, Violation{Key: key, Message: "required variable is missing or empty"})
				continue
			}
			continue
		}
		if err := validateType(key, val, field.Type); err != nil {
			violations = append(violations, Violation{Key: key, Message: err.Error()})
		}
	}
	return violations, nil
}

func validateType(key, val string, t FieldType) error {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("expected int, got %q", val)
		}
	case TypeBool:
		if _, err := strconv.ParseBool(val); err != nil {
			return fmt.Errorf("expected bool, got %q", val)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Errorf("expected float, got %q", val)
		}
	}
	return nil
}
