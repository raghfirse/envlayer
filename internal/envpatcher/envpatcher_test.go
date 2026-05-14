package envpatcher_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envpatcher"
)

func baseVars() map[string]string {
	return map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"DEBUG": "false",
	}
}

func TestApply_SetAddsNewKey(t *testing.T) {
	ops := []envpatcher.Op{{Kind: envpatcher.OpSet, Key: "TIMEOUT", Value: "30s"}}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %q", res.Vars["TIMEOUT"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(res.Applied))
	}
}

func TestApply_SetOverridesExistingKey(t *testing.T) {
	ops := []envpatcher.Op{{Kind: envpatcher.OpSet, Key: "PORT", Value: "3306"}}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["PORT"] != "3306" {
		t.Errorf("expected PORT=3306, got %q", res.Vars["PORT"])
	}
}

func TestApply_DeleteRemovesKey(t *testing.T) {
	ops := []envpatcher.Op{{Kind: envpatcher.OpDelete, Key: "DEBUG"}}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["DEBUG"]; ok {
		t.Error("expected DEBUG to be deleted")
	}
}

func TestApply_DeleteMissingKeyIsSkipped(t *testing.T) {
	ops := []envpatcher.Op{{Kind: envpatcher.OpDelete, Key: "NONEXISTENT"}}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped op, got %d", len(res.Skipped))
	}
}

func TestApply_RenameMovesKey(t *testing.T) {
	ops := []envpatcher.Op{{Kind: envpatcher.OpRename, Key: "HOST", NewKey: "DB_HOST"}}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", res.Vars["DB_HOST"])
	}
	if _, ok := res.Vars["HOST"]; ok {
		t.Error("expected original HOST key to be removed")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	orig := baseVars()
	ops := []envpatcher.Op{{Kind: envpatcher.OpSet, Key: "NEW", Value: "val"}}
	_, err := envpatcher.Apply(orig, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := orig["NEW"]; ok {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_UnknownOpReturnsError(t *testing.T) {
	ops := []envpatcher.Op{{Kind: "upsert", Key: "X"}}
	_, err := envpatcher.Apply(baseVars(), ops)
	if err == nil {
		t.Error("expected error for unknown op kind")
	}
}

func TestApply_MultipleOpsChained(t *testing.T) {
	ops := []envpatcher.Op{
		{Kind: envpatcher.OpSet, Key: "ENV", Value: "production"},
		{Kind: envpatcher.OpDelete, Key: "DEBUG"},
		{Kind: envpatcher.OpRename, Key: "PORT", NewKey: "DB_PORT"},
	}
	res, err := envpatcher.Apply(baseVars(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["ENV"] != "production" {
		t.Errorf("ENV mismatch")
	}
	if _, ok := res.Vars["DEBUG"]; ok {
		t.Error("DEBUG should be deleted")
	}
	if res.Vars["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT mismatch")
	}
	if len(res.Applied) != 3 {
		t.Errorf("expected 3 applied ops, got %d", len(res.Applied))
	}
}
