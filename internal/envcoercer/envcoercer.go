// Package envcoercer provides type coercion utilities for environment variable maps.
// It converts string values to target types based on a provided type hint map,
// returning a typed result map alongside any coercion errors encountered.
package envcoercer

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeHint describes the target type for a key.
type TypeHint string

const (
	TypeBool   TypeHint = "bool"
	TypeInt    TypeHint = "int"
	TypeFloat  TypeHint = "float"
	TypeString TypeHint = "string"
)

// CoercionError records a failure to coerce a single key.
type CoercionError struct {
	Key      string
	RawValue string
	Target   TypeHint
	Err      error
}

func (e CoercionError) Error() string {
	return fmt.Sprintf("coerce %q (%q -> %s): %v", e.Key, e.RawValue, e.Target, e.Err)
}

// Result holds the coerced values and any errors that occurred.
type Result struct {
	Values map[string]any
	Errors []CoercionError
}

// Coerce converts string values in vars according to hints.
// Keys without a hint are kept as strings.
// Keys with a hint that fail conversion are recorded in Result.Errors
// and the raw string value is retained in Result.Values.
func Coerce(vars map[string]string, hints map[string]TypeHint) Result {
	out := make(map[string]any, len(vars))
	var errs []CoercionError

	for k, v := range vars {
		hint, ok := hints[k]
		if !ok {
			out[k] = v
			continue
		}
		coerced, err := coerceValue(v, hint)
		if err != nil {
			errs = append(errs, CoercionError{Key: k, RawValue: v, Target: hint, Err: err})
			out[k] = v
			continue
		}
		out[k] = coerced
	}

	return Result{Values: out, Errors: errs}
}

func coerceValue(raw string, hint TypeHint) (any, error) {
	switch hint {
	case TypeBool:
		return strconv.ParseBool(strings.TrimSpace(raw))
	case TypeInt:
		return strconv.Atoi(strings.TrimSpace(raw))
	case TypeFloat:
		return strconv.ParseFloat(strings.TrimSpace(raw), 64)
	case TypeString:
		return raw, nil
	default:
		return nil, fmt.Errorf("unknown type hint %q", hint)
	}
}
