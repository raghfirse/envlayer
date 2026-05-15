// Package envboundary enforces key-level access boundaries between environment
// variable sets, allowing callers to define which keys are public, private, or
// restricted, and to assert that no boundary violations exist.
package envboundary

import (
	"fmt"
	"sort"
	"strings"
)

// Level represents the access level assigned to a key.
type Level int

const (
	Public     Level = iota // readable by anyone
	Internal                // readable within the same service boundary
	Restricted              // must not leak into downstream layers
)

// Violation describes a boundary rule that was broken.
type Violation struct {
	Key      string
	Level    Level
	Reason   string
}

func (v Violation) Error() string {
	return fmt.Sprintf("boundary violation: key %q (%s)", v.Key, v.Reason)
}

// Policy maps key prefixes or exact keys to their assigned Level.
type Policy struct {
	Rules map[string]Level // key or prefix -> level
}

// LevelFor returns the Level for a given key, matching exact keys first then
// the longest matching prefix. Defaults to Public if no rule matches.
func (p *Policy) LevelFor(key string) Level {
	if l, ok := p.Rules[key]; ok {
		return l
	}
	best, found := Public, false
	bestLen := 0
	for prefix, level := range p.Rules {
		if strings.HasPrefix(key, prefix) && len(prefix) > bestLen {
			best, found, bestLen = level, true, len(prefix)
		}
	}
	_ = found
	return best
}

// Check inspects vars against the policy and returns all Violations where a
// Restricted key appears in the provided map.
func Check(vars map[string]string, policy Policy, allowLevels ...Level) []Violation {
	allowed := map[Level]bool{}
	for _, l := range allowLevels {
		allowed[l] = true
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var violations []Violation
	for _, k := range keys {
		lvl := policy.LevelFor(k)
		if !allowed[lvl] {
			violations = append(violations, Violation{
				Key:    k,
				Level:  lvl,
				Reason: fmt.Sprintf("level %d not permitted in this context", lvl),
			})
		}
	}
	return violations
}

// Filter returns a copy of vars containing only keys whose Level is in the
// provided allow list.
func Filter(vars map[string]string, policy Policy, allowLevels ...Level) map[string]string {
	allowed := map[Level]bool{}
	for _, l := range allowLevels {
		allowed[l] = true
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if allowed[policy.LevelFor(k)] {
			out[k] = v
		}
	}
	return out
}
