package envscorer_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envscorer"
)

func TestScore_PerfectMap(t *testing.T) {
	vars := map[string]string{
		"APP_NAME": "myapp",
		"APP_ENV":  "production",
		"API_KEY":  "supersecret",
	}
	result := envscorer.Score(vars, envscorer.DefaultOptions())
	if result.Total != 100 {
		t.Errorf("expected 100, got %d; deductions: %v", result.Total, result.Deductions)
	}
}

func TestScore_EmptyMap_ReturnsZero(t *testing.T) {
	result := envscorer.Score(map[string]string{}, envscorer.DefaultOptions())
	if result.Total != 0 {
		t.Errorf("expected 0 for empty map, got %d", result.Total)
	}
	if len(result.Deductions) == 0 {
		t.Error("expected at least one deduction message")
	}
}

func TestScore_LowercaseKeyPenalty(t *testing.T) {
	vars := map[string]string{
		"app_name": "myapp",
	}
	opts := envscorer.DefaultOptions()
	opts.PenalizeLowercase = true
	result := envscorer.Score(vars, opts)
	if result.Total >= 100 {
		t.Errorf("expected penalty for lowercase key, got %d", result.Total)
	}
	found := false
	for _, d := range result.Deductions {
		if d == "key not uppercase: app_name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected deduction for lowercase key, got %v", result.Deductions)
	}
}

func TestScore_EmptyValuePenalty(t *testing.T) {
	vars := map[string]string{
		"APP_NAME": "",
		"APP_ENV":  "prod",
	}
	result := envscorer.Score(vars, envscorer.DefaultOptions())
	if result.Total >= 100 {
		t.Errorf("expected completeness penalty, got %d", result.Total)
	}
	if result.Categories["completeness"] >= 50 {
		t.Errorf("completeness category should be reduced, got %d", result.Categories["completeness"])
	}
}

func TestScore_SensitiveEmptyValuePenalty(t *testing.T) {
	vars := map[string]string{
		"APP_SECRET": "",
		"APP_NAME":   "myapp",
	}
	result := envscorer.Score(vars, envscorer.DefaultOptions())
	if result.Categories["security"] >= 25 {
		t.Errorf("expected security penalty, got %d", result.Categories["security"])
	}
	found := false
	for _, d := range result.Deductions {
		if d == "sensitive key has empty value: APP_SECRET" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected security deduction, got %v", result.Deductions)
	}
}

func TestScore_ShortKeyPenalty(t *testing.T) {
	vars := map[string]string{
		"AB": "value",
	}
	opts := envscorer.DefaultOptions()
	opts.PenalizeShortKeys = true
	opts.MinKeyLen = 3
	result := envscorer.Score(vars, opts)
	if result.Total >= 100 {
		t.Errorf("expected penalty for short key, got %d", result.Total)
	}
}

func TestScore_CategoriesSumCorrectly(t *testing.T) {
	vars := map[string]string{
		"APP_NAME": "myapp",
		"APP_ENV":  "prod",
	}
	result := envscorer.Score(vars, envscorer.DefaultOptions())
	sum := result.Categories["naming"] + result.Categories["completeness"] + result.Categories["security"]
	if sum != result.Total {
		t.Errorf("category sum %d does not match Total %d", sum, result.Total)
	}
}
