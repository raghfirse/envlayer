// Package envrender provides utilities for rendering environment variable
// maps into human-readable table or summary output formats.
package envrender

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls the output style of the rendered environment map.
type Format string

const (
	FormatTable   Format = "table"
	FormatCompact Format = "compact"
	FormatSummary Format = "summary"
)

// Options configures rendering behaviour.
type Options struct {
	Format      Format
	MaskValues  bool
	MaskChar    string
	MaxValueLen int
}

// DefaultOptions returns sensible rendering defaults.
func DefaultOptions() Options {
	return Options{
		Format:      FormatTable,
		MaskValues:  false,
		MaskChar:    "****",
		MaxValueLen: 64,
	}
}

// Render writes the env map to w using the given options.
func Render(w io.Writer, vars map[string]string, opts Options) error {
	keys := sortedKeys(vars)

	switch opts.Format {
	case FormatTable:
		return renderTable(w, keys, vars, opts)
	case FormatCompact:
		return renderCompact(w, keys, vars, opts)
	case FormatSummary:
		return renderSummary(w, keys, vars)
	default:
		return fmt.Errorf("envrender: unknown format %q", opts.Format)
	}
}

func renderTable(w io.Writer, keys []string, vars map[string]string, opts Options) error {
	const colPad = 2
	maxKey := 3 // minimum width "KEY"
	for _, k := range keys {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}
	fmt.Fprintf(w, "%-*s  %s\n", maxKey, "KEY", "VALUE")
	fmt.Fprintf(w, "%s  %s\n", strings.Repeat("-", maxKey), strings.Repeat("-", 20))
	for _, k := range keys {
		v := displayValue(vars[k], opts)
		fmt.Fprintf(w, "%-*s  %s\n", maxKey+colPad-2, k, v)
	}
	return nil
}

func renderCompact(w io.Writer, keys []string, vars map[string]string, opts Options) error {
	for _, k := range keys {
		v := displayValue(vars[k], opts)
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func renderSummary(w io.Writer, keys []string, vars map[string]string) error {
	fmt.Fprintf(w, "Total keys : %d\n", len(keys))
	empty := 0
	for _, k := range keys {
		if vars[k] == "" {
			empty++
		}
	}
	fmt.Fprintf(w, "Empty values: %d\n", empty)
	fmt.Fprintf(w, "Set values  : %d\n", len(keys)-empty)
	return nil
}

func displayValue(v string, opts Options) string {
	if opts.MaskValues {
		return opts.MaskChar
	}
	if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
		return v[:opts.MaxValueLen] + "..."
	}
	return v
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
