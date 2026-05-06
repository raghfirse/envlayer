// Package templater renders Go text/template strings using environment variables.
// It allows .env values to be used as template data, enabling dynamic config
// generation from the resolved environment context.
package templater

import (
	"bytes"
	"fmt"
	"text/template"
)

// Options configures template rendering behaviour.
type Options struct {
	// MissingKey controls behaviour when a key is absent from vars.
	// Accepted values: "zero" (default), "error", "invalid".
	MissingKey string
}

// DefaultOptions returns sensible defaults for template rendering.
func DefaultOptions() Options {
	return Options{MissingKey: "error"}
}

// Render parses tmpl as a Go text/template and executes it with vars as the
// data map. The rendered result is returned as a string.
//
// Example template: "host={{.DB_HOST}} port={{.DB_PORT}}"
func Render(tmpl string, vars map[string]string, opts Options) (string, error) {
	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}

	t, err := template.New("envlayer").
		Option("missingkey=" + missingKey).
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("templater: parse error: %w", err)
	}

	// Convert map[string]string to map[string]any so template engine
	// can access keys via {{.KEY}} notation.
	data := make(map[string]any, len(vars))
	for k, v := range vars {
		data[k] = v
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("templater: render error: %w", err)
	}

	return buf.String(), nil
}

// RenderFile renders the template string in tmpl using vars and writes the
// result to dest via the provided write function. This keeps I/O concerns
// outside the core rendering logic and simplifies testing.
func RenderFile(tmpl string, vars map[string]string, opts Options, write func(string) error) error {
	out, err := Render(tmpl, vars, opts)
	if err != nil {
		return err
	}
	return write(out)
}
