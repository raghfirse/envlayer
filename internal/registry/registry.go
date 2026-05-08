package registry

import (
	"fmt"
	"sort"
	"sync"
)

// Entry holds a named environment map with optional metadata tags.
type Entry struct {
	Name string
	Vars map[string]string
	Tags []string
}

// Registry is an in-memory store of named environment maps.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[string]*Entry)}
}

// Register stores an environment map under the given name, replacing any
// existing entry with the same name.
func (r *Registry) Register(name string, vars map[string]string, tags ...string) {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[name] = &Entry{Name: name, Vars: copy, Tags: tags}
}

// Get retrieves an entry by name. Returns an error if not found.
func (r *Registry) Get(name string) (*Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[name]
	if !ok {
		return nil, fmt.Errorf("registry: entry %q not found", name)
	}
	return e, nil
}

// Remove deletes an entry by name. Returns an error if not found.
func (r *Registry) Remove(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[name]; !ok {
		return fmt.Errorf("registry: entry %q not found", name)
	}
	delete(r.entries, name)
	return nil
}

// Names returns a sorted list of all registered entry names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.entries))
	for k := range r.entries {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// FindByTag returns all entries that contain the given tag, sorted by name.
func (r *Registry) FindByTag(tag string) []*Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Entry
	for _, e := range r.entries {
		for _, t := range e.Tags {
			if t == tag {
				result = append(result, e)
				break
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
