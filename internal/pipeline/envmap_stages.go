package pipeline

import (
	"strings"

	"github.com/yourorg/envlayer/internal/envmap"
)

// StageFilterKeys returns a Stage that keeps only entries whose keys satisfy
// the given predicate, using envmap.Filter for ordered processing.
func StageFilterKeys(predicate func(key string) bool) Stage {
	return Stage{
		Name: "filter_keys",
		Run: func(vars map[string]string) (map[string]string, error) {
			entries := envmap.Filter(envmap.FromMap(vars), predicate)
			return envmap.ToMap(entries), nil
		},
	}
}

// StageTransformValues returns a Stage that applies fn to every value.
func StageTransformValues(name string, fn func(key, value string) string) Stage {
	return Stage{
		Name: name,
		Run: func(vars map[string]string) (map[string]string, error) {
			entries := envmap.MapValues(envmap.FromMap(vars), fn)
			return envmap.ToMap(entries), nil
		},
	}
}

// StageTrimValues returns a Stage that trims leading/trailing whitespace from
// all values using the envmap transform pipeline.
func StageTrimValues() Stage {
	return StageTransformValues("trim_values", func(_, v string) string {
		return strings.TrimSpace(v)
	})
}
