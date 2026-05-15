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
package envproxy
