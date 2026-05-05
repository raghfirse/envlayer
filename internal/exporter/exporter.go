// Package exporter provides functionality for exporting merged environment
// variables in various output formats (shell export statements, dotenv, JSON).
package exporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	// FormatExport produces shell-compatible `export KEY=VALUE` lines.
	FormatExport Format = "export"
	// FormatDotenv produces standard KEY=VALUE lines suitable for a .env file.
	FormatDotenv Format = "dotenv"
	// FormatJSON produces a JSON object of key/value pairs.
	FormatJSON Format = "json"
)

// Export serialises the given environment map into the requested format.
// Keys are always emitted in sorted order for deterministic output.
func Export(env map[string]string, format Format) (string, error) {
	switch format {
	case FormatExport:
		return exportShell(env), nil
	case FormatDotenv:
		return exportDotenv(env), nil
	case FormatJSON:
		return exportJSON(env)
	default:
		return "", fmt.Errorf("exporter: unknown format %q", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func exportShell(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportDotenv(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportJSON(env map[string]string) (string, error) {
	// Use a sorted intermediate structure so output is deterministic.
	ordered := make(map[string]string, len(env))
	for k, v := range env {
		ordered[k] = v
	}
	b, err := json.MarshalIndent(ordered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("exporter: json marshal failed: %w", err)
	}
	return string(b) + "\n", nil
}
