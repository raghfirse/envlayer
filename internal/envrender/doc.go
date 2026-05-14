// Package envrender renders environment variable maps into human-readable
// output formats for display in terminals or log output.
//
// Three formats are supported:
//
//   - table   – aligned two-column table with KEY / VALUE header (default)
//   - compact – KEY=VALUE lines, one per entry
//   - summary – aggregate statistics: total keys, empty values, set values
//
// Values can be masked (e.g. for secrets) and truncated to a maximum length
// to keep output readable. All formats write to any io.Writer, making them
// easy to integrate with CLI commands, HTTP handlers, or log sinks.
//
// Example:
//
//	opts := envrender.DefaultOptions()
//	opts.MaskValues = true
//	_ = envrender.Render(os.Stdout, vars, opts)
package envrender
