package envcoercer_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envcoercer"
)

func TestCoerce_StringPassthrough(t *testing.T) {
	vars := map[string]string{"NAME": "alice"}
	res := envcoercer.Coerce(vars, map[string]envcoercer.TypeHint{})
	if v, ok := res.Values["NAME"]; !ok || v != "alice" {
		t.Fatalf("expected 'alice', got %v", v)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestCoerce_IntSuccess(t *testing.T) {
	vars := map[string]string{"PORT": "8080"}
	hints := map[string]envcoercer.TypeHint{"PORT": envcoercer.TypeInt}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Values["PORT"] != 8080 {
		t.Fatalf("expected 8080, got %v", res.Values["PORT"])
	}
}

func TestCoerce_BoolSuccess(t *testing.T) {
	vars := map[string]string{"DEBUG": "true"}
	hints := map[string]envcoercer.TypeHint{"DEBUG": envcoercer.TypeBool}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Values["DEBUG"] != true {
		t.Fatalf("expected true, got %v", res.Values["DEBUG"])
	}
}

func TestCoerce_FloatSuccess(t *testing.T) {
	vars := map[string]string{"RATIO": "3.14"}
	hints := map[string]envcoercer.TypeHint{"RATIO": envcoercer.TypeFloat}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	v, ok := res.Values["RATIO"].(float64)
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %v", res.Values["RATIO"])
	}
}

func TestCoerce_InvalidInt_RecordsError(t *testing.T) {
	vars := map[string]string{"PORT": "not-a-number"}
	hints := map[string]envcoercer.TypeHint{"PORT": envcoercer.TypeInt}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(res.Errors))
	}
	if res.Errors[0].Key != "PORT" {
		t.Fatalf("expected error on PORT, got %q", res.Errors[0].Key)
	}
	// raw value retained
	if res.Values["PORT"] != "not-a-number" {
		t.Fatalf("expected raw value retained, got %v", res.Values["PORT"])
	}
}

func TestCoerce_InvalidBool_RecordsError(t *testing.T) {
	vars := map[string]string{"FLAG": "maybe"}
	hints := map[string]envcoercer.TypeHint{"FLAG": envcoercer.TypeBool}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestCoerce_MixedHints(t *testing.T) {
	vars := map[string]string{"PORT": "9000", "NAME": "svc", "DEBUG": "false"}
	hints := map[string]envcoercer.TypeHint{
		"PORT":  envcoercer.TypeInt,
		"DEBUG": envcoercer.TypeBool,
	}
	res := envcoercer.Coerce(vars, hints)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Values["PORT"] != 9000 {
		t.Fatalf("expected 9000, got %v", res.Values["PORT"])
	}
	if res.Values["DEBUG"] != false {
		t.Fatalf("expected false, got %v", res.Values["DEBUG"])
	}
	if res.Values["NAME"] != "svc" {
		t.Fatalf("expected 'svc', got %v", res.Values["NAME"])
	}
}
