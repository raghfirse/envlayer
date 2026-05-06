// Package snapshot provides functionality to capture and persist
// the current state of merged environment variables to a file.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Snapshot represents a captured state of environment variables.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Environment string          `json:"environment,omitempty"`
	Vars      map[string]string `json:"vars"`
}

// Take creates a new Snapshot from the given environment map.
func Take(env map[string]string, environment string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		CreatedAt:   time.Now().UTC(),
		Environment: environment,
		Vars:        copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads and parses a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: parse failed: %w", err)
	}
	return &s, nil
}

// SortedKeys returns the variable keys in sorted order.
func (s *Snapshot) SortedKeys() []string {
	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
