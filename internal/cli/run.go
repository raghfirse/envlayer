// Package cli provides the command-line interface for envlayer.
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/envlayer/envlayer/internal/exporter"
	"github.com/envlayer/envlayer/internal/loader"
	"github.com/envlayer/envlayer/internal/printer"
	"github.com/envlayer/envlayer/internal/resolver"
	"github.com/envlayer/envlayer/internal/validator"
)

// Config holds the parsed CLI flags.
type Config struct {
	Env      string
	Dir      string
	Format   string
	Prefix   string
	Required string
	Export   bool
}

// Run parses flags and executes the envlayer pipeline.
func Run(args []string) int {
	fs := flag.NewFlagSet("envlayer", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	cfg := &Config{}
	fs.StringVar(&cfg.Env, "env", "", "environment name (e.g. staging, production)")
	fs.StringVar(&cfg.Dir, "dir", ".", "directory to search for .env files")
	fs.StringVar(&cfg.Format, "format", "dotenv", "output format: shell, dotenv, json")
	fs.StringVar(&cfg.Prefix, "prefix", "", "filter keys by prefix")
	fs.StringVar(&cfg.Required, "require", "", "comma-separated list of required keys")
	fs.BoolVar(&cfg.Export, "export", false, "write to stdout using export format")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	files, err := resolver.Resolve(cfg.Dir, cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envlayer: resolve error: %v\n", err)
		return 1
	}

	env, err := loader.LoadFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envlayer: load error: %v\n", err)
		return 1
	}

	if cfg.Required != "" {
		requiredKeys := splitCSV(cfg.Required)
		result := validator.Validate(env, requiredKeys)
		for _, w := range result.Warnings {
			fmt.Fprintf(os.Stderr, "envlayer: warning: %s\n", w)
		}
		if len(result.Missing) > 0 {
			for _, m := range result.Missing {
				fmt.Fprintf(os.Stderr, "envlayer: missing required key: %s\n", m)
			}
			return 1
		}
	}

	if cfg.Export {
		if err := exporter.Export(os.Stdout, env, exporter.Format(cfg.Format)); err != nil {
			fmt.Fprintf(os.Stderr, "envlayer: export error: %v\n", err)
			return 1
		}
		return 0
	}

	p := printer.New(os.Stdout)
	if err := p.Print(env, cfg.Format, cfg.Prefix); err != nil {
		fmt.Fprintf(os.Stderr, "envlayer: print error: %v\n", err)
		return 1
	}
	return 0
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			if tok := s[start:i]; tok != "" {
				out = append(out, tok)
			}
			start = i + 1
		}
	}
	if tok := s[start:]; tok != "" {
		out = append(out, tok)
	}
	return out
}
