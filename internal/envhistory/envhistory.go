package envhistory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single historical snapshot of environment variables.
type Entry struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Record appends a new history entry to the given directory.
func Record(dir, label string, vars map[string]string) (Entry, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Entry{}, fmt.Errorf("envhistory: mkdir: %w", err)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	e := Entry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Label:     label,
		Vars:      copy,
		CreatedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return Entry{}, fmt.Errorf("envhistory: marshal: %w", err)
	}

	path := filepath.Join(dir, e.ID+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return Entry{}, fmt.Errorf("envhistory: write: %w", err)
	}
	return e, nil
}

// List returns all history entries in the directory sorted by creation time (oldest first).
func List(dir string) ([]Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("envhistory: readdir: %w", err)
	}

	var result []Entry
	for _, de := range entries {
		if de.IsDir() || filepath.Ext(de.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, de.Name()))
		if err != nil {
			return nil, fmt.Errorf("envhistory: read %s: %w", de.Name(), err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("envhistory: unmarshal %s: %w", de.Name(), err)
		}
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

// Get retrieves a single entry by ID.
func Get(dir, id string) (Entry, error) {
	path := filepath.Join(dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Entry{}, fmt.Errorf("envhistory: entry %q not found", id)
		}
		return Entry{}, fmt.Errorf("envhistory: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, fmt.Errorf("envhistory: unmarshal: %w", err)
	}
	return e, nil
}

// Delete removes a history entry by ID.
func Delete(dir, id string) error {
	path := filepath.Join(dir, id+".json")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("envhistory: entry %q not found", id)
		}
		return fmt.Errorf("envhistory: delete: %w", err)
	}
	return nil
}
