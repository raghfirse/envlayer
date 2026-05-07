// Package cascader provides hierarchical .env file resolution and merging
// across multiple named environments with inheritance support.
package cascader

import (
	"fmt"

	"github.com/yourusername/envlayer/internal/loader"
	"github.com/yourusername/envlayer/internal/merger"
)

// Layer represents a named environment layer with its loaded variables.
type Layer struct {
	Name string
	Vars map[string]string
}

// CascadeOptions controls how the cascade is built.
type CascadeOptions struct {
	// Dir is the directory to search for .env files.
	Dir string
	// Environments is the ordered list of environment names to cascade through.
	// Variables from later environments override earlier ones.
	Environments []string
}

// Result holds the final merged variables and the ordered layers used.
type Result struct {
	Vars   map[string]string
	Layers []Layer
}

// Build resolves and merges .env files for each environment in order.
// For each environment name, it attempts to load "<dir>/.env.<name>".
// The base ".env" file is always loaded first if present.
func Build(opts CascadeOptions) (*Result, error) {
	if opts.Dir == "" {
		return nil, fmt.Errorf("cascader: Dir must not be empty")
	}

	var layers []Layer
	var maps []map[string]string

	// Always attempt to load the base .env first.
	baseFile := opts.Dir + "/.env"
	baseVars, err := loader.LoadFile(baseFile)
	if err == nil {
		layers = append(layers, Layer{Name: "base", Vars: baseVars})
		maps = append(maps, baseVars)
	}

	for _, env := range opts.Environments {
		if env == "" {
			continue
		}
		path := fmt.Sprintf("%s/.env.%s", opts.Dir, env)
		vars, err := loader.LoadFile(path)
		if err != nil {
			// Non-fatal: skip missing layers.
			continue
		}
		layers = append(layers, Layer{Name: env, Vars: vars})
		maps = append(maps, vars)
	}

	if len(maps) == 0 {
		return nil, fmt.Errorf("cascader: no .env files found in %q for environments %v", opts.Dir, opts.Environments)
	}

	merged := merger.Merge(maps...)
	return &Result{Vars: merged, Layers: layers}, nil
}
