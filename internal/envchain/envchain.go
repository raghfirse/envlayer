// Package envchain provides a chainable, ordered collection of environment
// variable maps that resolves keys by searching layers from highest to lowest
// priority, similar to a prototype chain.
package envchain

import "fmt"

// Chain holds an ordered list of env layers where index 0 is highest priority.
type Chain struct {
	layers []map[string]string
}

// New creates a new Chain from the provided layers. The first layer has the
// highest priority; subsequent layers act as fallbacks.
func New(layers ...map[string]string) *Chain {
	copy := make([]map[string]string, len(layers))
	for i, l := range layers {
		copy[i] = cloneMap(l)
	}
	return &Chain{layers: copy}
}

// Get returns the value for key by searching layers in priority order.
// Returns an error if the key is not found in any layer.
func (c *Chain) Get(key string) (string, error) {
	for _, layer := range c.layers {
		if v, ok := layer[key]; ok {
			return v, nil
		}
	}
	return "", fmt.Errorf("envchain: key %q not found in any layer", key)
}

// GetOrDefault returns the value for key or fallback if not found.
func (c *Chain) GetOrDefault(key, fallback string) string {
	if v, err := c.Get(key); err == nil {
		return v
	}
	return fallback
}

// Resolve collapses all layers into a single map, with higher-priority layers
// winning on key conflicts.
func (c *Chain) Resolve() map[string]string {
	out := make(map[string]string)
	for i := len(c.layers) - 1; i >= 0; i-- {
		for k, v := range c.layers[i] {
			out[k] = v
		}
	}
	return out
}

// Len returns the number of layers in the chain.
func (c *Chain) Len() int { return len(c.layers) }

// Push adds a new highest-priority layer to the chain.
func (c *Chain) Push(layer map[string]string) {
	c.layers = append([]map[string]string{cloneMap(layer)}, c.layers...)
}

func cloneMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
