package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envlayer/internal/exporter"
	"github.com/yourusername/envlayer/internal/loader"
	"github.com/yourusername/envlayer/internal/masker"
	"github.com/yourusername/envlayer/internal/resolver"
	"github.com/yourusername/envlayer/internal/validator"
)

// Config holds all CLI options passed to Run.
type Config struct {
	Dir         string
	Environment string
	Format      string
	Required    string
	Prefix      string
	MaskOutput  bool
	MaskKeys    string // comma-separated extra sensitive key patterns
}

// Run is the main entry point for the envlayer CLI.
func Run(cfg Config, out *os.File) error {
	files, err := resolver.Resolve(cfg.Dir, cfg.Environment)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	vars, err := loader.LoadFiles(files)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if cfg.Required != "" {
		keys := splitCSV(cfg.Required)
		result := validator.Validate(vars, keys)
		if len(result.Missing) > 0 {
			return fmt.Errorf("missing required keys: %s", strings.Join(result.Missing, ", "))
		}
	}

	output := vars
	if cfg.MaskOutput {
		opts := masker.Options{}
		if cfg.MaskKeys != "" {
			opts.SensitiveKeys = append(masker.DefaultSensitiveKeys, splitCSV(cfg.MaskKeys)...)
		}
		output = masker.Mask(vars, opts)
	}

	fmt := cfg.Format
	if fmt == "" {
		fmt = "dotenv"
	}

	return exporter.Export(output, exporter.Options{
		Format: fmt,
		Prefix: cfg.Prefix,
	}, out)
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
