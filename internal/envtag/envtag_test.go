package envtag_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/envtag"
)

var sampleVars = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"API_KEY":     "abc123",
	"LOG_LEVEL":   "info",
}

var sampleTags = []envtag.Tag{
	{Name: "database", Keys: []string{"DB_HOST", "DB_PASSWORD"}},
	{Name: "secret", Keys: []string{"DB_PASSWORD", "API_KEY"}},
	{Name: "logging", Keys: []string{"LOG_LEVEL"}},
}

func TestBuild_IndexesKeysByTag(t *testing.T) {
	idx := envtag.Build(sampleTags)
	tags := envtag.TagsFor(idx, "DB_PASSWORD")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags for DB_PASSWORD, got %d", len(tags))
	}
	if tags[0] != "database" || tags[1] != "secret" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestFilterByTag_ReturnsMatchingKeys(t *testing.T) {
	idx := envtag.Build(sampleTags)
	result := envtag.FilterByTag(sampleVars, idx, "database")
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", result["DB_HOST"])
	}
	if result["DB_PASSWORD"] != "secret" {
		t.Errorf("unexpected value for DB_PASSWORD: %s", result["DB_PASSWORD"])
	}
}

func TestFilterByTag_UnknownTagReturnsEmpty(t *testing.T) {
	idx := envtag.Build(sampleTags)
	result := envtag.FilterByTag(sampleVars, idx, "nonexistent")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestAllTags_ReturnsSortedUnique(t *testing.T) {
	idx := envtag.Build(sampleTags)
	all := envtag.AllTags(idx)
	expected := []string{"database", "logging", "secret"}
	if len(all) != len(expected) {
		t.Fatalf("expected %d tags, got %d: %v", len(expected), len(all), all)
	}
	for i, tag := range expected {
		if all[i] != tag {
			t.Errorf("position %d: expected %q, got %q", i, tag, all[i])
		}
	}
}

func TestValidate_AllKeysPresent(t *testing.T) {
	idx := envtag.Build(sampleTags)
	if err := envtag.Validate(sampleVars, idx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingKeyReturnsError(t *testing.T) {
	idx := envtag.Build([]envtag.Tag{
		{Name: "extra", Keys: []string{"MISSING_KEY"}},
	})
	err := envtag.Validate(sampleVars, idx)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTagsFor_UnknownKeyReturnsEmpty(t *testing.T) {
	idx := envtag.Build(sampleTags)
	tags := envtag.TagsFor(idx, "UNKNOWN")
	if len(tags) != 0 {
		t.Errorf("expected empty slice, got %v", tags)
	}
}
