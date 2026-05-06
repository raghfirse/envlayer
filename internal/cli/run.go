package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/user/envlayer/internal/exporter"
	"github.com/user/envlayer/internal/interpolator"
	"github.com/user/envlayer/internal/loader"
	"github.com/user/envlayer/internal/printer"
	"github.com/user/envlayer/internal/resolver"
	"github.com/user/envlayer/internal/validator"
)

// Config holds the parsed CLI options passed to Run.
type Config struct {
	Dir         string
	Environment string
	Format      string
	Prefix      string
	Required    []string
	Interpolate bool
	Output      *os.File
}

// Run executes the core envlayer pipeline:
// resolve → load → (interpolate) → validate → print.
func Run(cfg Config) error {
	files, err := resolver.Resolve(cfg.Dir, cfg.Environment)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	env, err := loader.LoadFiles(files)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if cfg.Interpolate {
		env, err = interpolator.Interpolate(env, nil, interpolator.Options{
			FallbackToOS: true,
		})
		if err != nil {
			return fmt.Errorf("interpolate: %w", err)
		}
	}

	if len(cfg.Required) > 0 {
		result := validator.Validate(env, cfg.Required)
		for _, w := range result.Warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
		if len(result.Missing) > 0 {
			return fmt.Errorf("missing required keys: %s", strings.Join(result.Missing, ", "))
		}
	}

	out := cfg.Output
	if out == nil {
		out = os.Stdout
	}

	p := printer.New(out)
	return p.Print(env, printer.Options{
		Format: cfg.Format,
		Prefix: cfg.Prefix,
	})
}

// splitCSV splits a comma-separated string into a trimmed slice.
// Empty input returns nil.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// errNotFound is a sentinel used in tests.
var errNotFound = errors.New("not found")
