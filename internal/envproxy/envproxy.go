// Package envproxy provides a read-through proxy layer over an environment
// variable map, optionally delegating missing key lookups to a fallback source
// such as os.Environ or another map.
package envproxy

import (
	"fmt"
	"os"
	"sort"
)

// FallbackFunc is called when a key is not found in the primary map.
type FallbackFunc func(key string) (string, bool)

// Proxy wraps a primary env map and delegates missing lookups to a fallback.
type Proxy struct {
	primary  map[string]string
	fallback FallbackFunc
}

// New creates a Proxy backed by primary. If fallback is nil, missing keys
// return an error without further lookup.
func New(primary map[string]string, fallback FallbackFunc) *Proxy {
	copy := make(map[string]string, len(primary))
	for k, v := range primary {
		copy[k] = v
	}
	return &Proxy{primary: copy, fallback: fallback}
}

// WithOSFallback returns a Proxy that falls back to os.LookupEnv.
func WithOSFallback(primary map[string]string) *Proxy {
	return New(primary, os.LookupEnv)
}

// Get returns the value for key. It checks the primary map first, then the
// fallback. Returns an error if the key is not found in either source.
func (p *Proxy) Get(key string) (string, error) {
	if v, ok := p.primary[key]; ok {
		return v, nil
	}
	if p.fallback != nil {
		if v, ok := p.fallback(key); ok {
			return v, nil
		}
	}
	return "", fmt.Errorf("envproxy: key %q not found", key)
}

// GetOrDefault returns the value for key, or def if not found.
func (p *Proxy) GetOrDefault(key, def string) string {
	v, err := p.Get(key)
	if err != nil {
		return def
	}
	return v
}

// Has reports whether key is available (primary or fallback).
func (p *Proxy) Has(key string) bool {
	_, err := p.Get(key)
	return err == nil
}

// Keys returns the sorted keys of the primary map only.
func (p *Proxy) Keys() []string {
	keys := make([]string, 0, len(p.primary))
	for k := range p.primary {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Resolve returns a snapshot map of all primary keys, with missing values
// resolved through the fallback where possible.
func (p *Proxy) Resolve() map[string]string {
	out := make(map[string]string, len(p.primary))
	for k, v := range p.primary {
		out[k] = v
	}
	return out
}
