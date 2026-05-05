package validator_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlayer/internal/validator"
)

func TestValidate_AllRequiredPresent(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_URL":  "postgres://localhost/mydb",
	}
	res := validator.Validate(env, []string{"APP_ENV", "DB_URL"})
	if !res.IsValid() {
		t.Fatalf("expected valid, got error: %s", res.Error())
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "staging",
	}
	res := validator.Validate(env, []string{"APP_ENV", "DB_URL"})
	if res.IsValid() {
		t.Fatal("expected invalid result")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_URL" {
		t.Fatalf("expected Missing=[DB_URL], got %v", res.Missing)
	}
}

func TestValidate_EmptyValueCountsAsMissing(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "",
		"DB_URL":  "postgres://localhost/mydb",
	}
	res := validator.Validate(env, []string{"APP_ENV", "DB_URL"})
	if res.IsValid() {
		t.Fatal("expected invalid because APP_ENV is empty")
	}
	if !strings.Contains(res.Error(), "APP_ENV") {
		t.Fatalf("error should mention APP_ENV, got: %s", res.Error())
	}
}

func TestValidate_NoRequiredKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	res := validator.Validate(env, nil)
	if !res.IsValid() {
		t.Fatalf("expected valid with no required keys, got: %s", res.Error())
	}
}

func TestValidate_EmptyOptionalKeyGeneratesWarning(t *testing.T) {
	env := map[string]string{
		"REQUIRED_KEY": "value",
		"OPTIONAL_KEY": "",
	}
	res := validator.Validate(env, []string{"REQUIRED_KEY"})
	if !res.IsValid() {
		t.Fatalf("expected valid result, got: %s", res.Error())
	}
	if len(res.Warnings) == 0 {
		t.Fatal("expected a warning for empty optional key")
	}
	if !strings.Contains(res.Warnings[0], "OPTIONAL_KEY") {
		t.Fatalf("warning should mention OPTIONAL_KEY, got: %v", res.Warnings)
	}
}

func TestResult_ErrorIsEmptyWhenValid(t *testing.T) {
	res := validator.Result{}
	if res.Error() != "" {
		t.Fatalf("expected empty error string, got: %q", res.Error())
	}
}
