package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlayer/internal/pipeline"
)

func TestStageFilterKeys_KeepsMatchingKeys(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "secret",
		"APP_PORT": "9000",
	}
	p := pipeline.New(pipeline.StageFilterKeys(func(k string) bool {
		return strings.HasPrefix(k, "APP_")
	}))
	result, err := p.Run(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_PASS"]; ok {
		t.Error("DB_PASS should have been filtered out")
	}
}

func TestStageTransformValues_AppliesFn(t *testing.T) {
	vars := map[string]string{"MSG": "hello", "NAME": "world"}
	p := pipeline.New(pipeline.StageTransformValues("upper_values", func(_, v string) string {
		return strings.ToUpper(v)
	}))
	result, err := p.Run(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["MSG"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", result["MSG"])
	}
	if result["NAME"] != "WORLD" {
		t.Errorf("expected WORLD, got %q", result["NAME"])
	}
}

func TestStageTrimValues_TrimsWhitespace(t *testing.T) {
	vars := map[string]string{
		"HOST": "  localhost  ",
		"PORT": "\t8080\n",
		"CLEAN": "ok",
	}
	p := pipeline.New(pipeline.StageTrimValues())
	result, err := p.Run(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", result["PORT"])
	}
	if result["CLEAN"] != "ok" {
		t.Errorf("expected 'ok', got %q", result["CLEAN"])
	}
}

func TestStageFilterKeys_EmptyMap(t *testing.T) {
	p := pipeline.New(pipeline.StageFilterKeys(func(k string) bool { return true }))
	result, err := p.Run(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
