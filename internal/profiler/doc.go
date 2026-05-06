// Package profiler provides named environment profile management for envlayer.
//
// A profile is a named snapshot of environment variables that can be saved to
// disk, loaded by name, listed, and deleted. Profiles are stored as JSON files
// in a configurable directory (e.g. ~/.config/envlayer/profiles/).
//
// Typical usage:
//
//	// Save the current resolved vars as a profile named "staging".
//	_ = profiler.Save(profilesDir, "staging", vars)
//
//	// Load a previously saved profile.
//	p, _ := profiler.Load(profilesDir, "staging")
//
//	// List all available profiles.
//	names, _ := profiler.List(profilesDir)
//
//	// Remove a profile.
//	_ = profiler.Delete(profilesDir, "staging")
package profiler
