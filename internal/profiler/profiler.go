// Package profiler provides environment profile management, allowing
// named profiles to be saved, loaded, and switched between environments.
package profiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Profile represents a named set of environment variables.
type Profile struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Save writes a profile to a JSON file in the given directory.
func Save(dir, name string, vars map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("profiler: mkdir %s: %w", dir, err)
	}

	p := Profile{
		Name:      name,
		Vars:      vars,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("profiler: marshal: %w", err)
	}

	path := filepath.Join(dir, name+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("profiler: write %s: %w", path, err)
	}
	return nil
}

// Load reads a named profile from the given directory.
func Load(dir, name string) (*Profile, error) {
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profiler: profile %q not found", name)
		}
		return nil, fmt.Errorf("profiler: read %s: %w", path, err)
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("profiler: unmarshal: %w", err)
	}
	return &p, nil
}

// List returns the names of all saved profiles in the given directory.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("profiler: readdir %s: %w", dir, err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	sort.Strings(names)
	return names, nil
}

// Delete removes a named profile from the given directory.
func Delete(dir, name string) error {
	path := filepath.Join(dir, name+".json")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("profiler: profile %q not found", name)
		}
		return fmt.Errorf("profiler: delete %s: %w", path, err)
	}
	return nil
}
