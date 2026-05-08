// Package registry provides an in-memory, concurrency-safe store for named
// environment maps.
//
// A Registry allows multiple resolved environment snapshots to be stored,
// retrieved, and queried by name or tag. This is useful when an application
// needs to manage several environment contexts simultaneously — for example,
// keeping "dev", "staging", and "prod" configurations loaded at once and
// switching between them without re-reading files from disk.
//
// Usage:
//
//	reg := registry.New()
//	reg.Register("prod", vars, "live", "stable")
//
//	entry, err := reg.Get("prod")
//	live := reg.FindByTag("live")
//	names := reg.Names()
package registry
