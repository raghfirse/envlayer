// Package envproxy implements a read-through proxy for environment variable
// maps.
//
// A Proxy wraps a primary map[string]string and, on a cache miss, delegates
// the lookup to a configurable FallbackFunc. This allows callers to layer
// multiple sources — e.g. a loaded .env file backed by the real OS environment
// — without merging them upfront.
//
// Usage:
//
//	p := envproxy.WithOSFallback(loaded)
//	v, err := p.Get("DATABASE_URL")
//
// The primary map is copied on construction so mutations to the original map
// do not affect the proxy.
//
// Lookup semantics:
//
//  1. If the key exists in the primary map, its value is returned immediately
//     (even if the value is an empty string).
//  2. Otherwise, FallbackFunc is called with the key. If FallbackFunc is nil,
//     Get returns ("", ErrNotFound).
//  3. If FallbackFunc returns ok=false, Get returns ("", ErrNotFound).
package envproxy
