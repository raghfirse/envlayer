// Package printer handles formatted output of merged environment variables
// to various destinations such as stdout or a file.
package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envlayer/internal/exporter"
)

// Options configures the behavior of the Printer.
type Options struct {
	// Format is the output format: "shell", "dotenv", or "json".
	Format string
	// Output is the writer to write to. Defaults to os.Stdout if nil.
	Output io.Writer
	// Prefix optionally filters keys to only those with the given prefix.
	Prefix string
}

// Printer writes exported environment variables to an output destination.
type Printer struct {
	opts Options
}

// New creates a new Printer with the given options.
func New(opts Options) *Printer {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.Format == "" {
		opts.Format = "shell"
	}
	return &Printer{opts: opts}
}

// Print filters and writes the environment map using the configured format.
func (p *Printer) Print(env map[string]string) error {
	filtered := env
	if p.opts.Prefix != "" {
		filtered = filterByPrefix(env, p.opts.Prefix)
	}

	out, err := exporter.Export(filtered, p.opts.Format)
	if err != nil {
		return fmt.Errorf("printer: export failed: %w", err)
	}

	_, err = fmt.Fprint(p.opts.Output, out)
	if err != nil {
		return fmt.Errorf("printer: write failed: %w", err)
	}
	return nil
}

// filterByPrefix returns a new map containing only keys that start with prefix.
func filterByPrefix(env map[string]string, prefix string) map[string]string {
	result := make(map[string]string)
	for k, v := range env {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			result[k] = v
		}
	}
	return result
}
