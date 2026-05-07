package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlayer/envlayer/internal/pipeline"
	"github.com/envlayer/envlayer/internal/resolver"
	"github.com/envlayer/envlayer/internal/loader"
	"github.com/envlayer/envlayer/internal/printer"
)

// PipelineOptions configures the pipeline command execution.
type PipelineOptions struct {
	// Dir is the directory containing .env files.
	Dir string
	// Environment selects the environment layer (e.g. "production").
	Environment string
	// Stages is a comma-separated list of pipeline stage names to apply.
	// Supported: interpolate, mask, prefix:<value>, uppercase, filter:<prefix>
	Stages string
	// Format controls the output format: shell, dotenv, or json.
	Format string
	// Prefix is used by the prefix and filter stages.
	Prefix string
	// Output is the writer for results; defaults to os.Stdout.
	Output *os.File
}

// RunPipeline resolves env files, loads them, and runs the requested pipeline
// stages before printing the resulting variable map.
func RunPipeline(opts PipelineOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	// Resolve candidate files for the given directory and environment.
	files, err := resolver.Resolve(opts.Dir, opts.Environment)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	// Load and merge all resolved files.
	vars, err := loader.LoadFiles(files)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	// Build the pipeline from the requested stage names.
	p, err := buildPipeline(opts.Stages, opts.Prefix)
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	// Execute the pipeline.
	result, err := p.Run(vars)
	if err != nil {
		return fmt.Errorf("pipeline run: %w", err)
	}

	// Print the final variable map.
	pr := printer.New(opts.Output)
	if err := pr.Print(result, opts.Format, ""); err != nil {
		return fmt.Errorf("print: %w", err)
	}

	return nil
}

// buildPipeline constructs a pipeline.Pipeline from a comma-separated list of
// stage descriptors. Each descriptor is either a plain name (e.g. "mask") or a
// name with an argument separated by a colon (e.g. "prefix:APP_").
func buildPipeline(stages string, defaultPrefix string) (*pipeline.Pipeline, error) {
	p := pipeline.New()

	if strings.TrimSpace(stages) == "" {
		return p, nil
	}

	for _, raw := range splitCSV(stages) {
		name, arg, _ := strings.Cut(raw, ":")
		name = strings.TrimSpace(name)
		arg = strings.TrimSpace(arg)

		switch name {
		case "interpolate":
			p.Add("interpolate", pipeline.StageInterpolate)
		case "mask":
			p.Add("mask", pipeline.StageMask)
		case "uppercase":
			p.Add("uppercase", pipeline.StageUppercaseKeys)
		case "trim":
			p.Add("trim", pipeline.StageTrimValues)
		case "prefix":
			pfx := arg
			if pfx == "" {
				pfx = defaultPrefix
			}
			if pfx == "" {
				return nil, fmt.Errorf("stage %q requires a prefix argument (e.g. prefix:APP_)", name)
			}
			p.Add("prefix", pipeline.StageAddPrefix(pfx))
		case "filter":
			pfx := arg
			if pfx == "" {
				pfx = defaultPrefix
			}
			if pfx == "" {
				return nil, fmt.Errorf("stage %q requires a prefix argument (e.g. filter:APP_)", name)
			}
			p.Add("filter", pipeline.StageFilterPrefix(pfx))
		default:
			return nil, fmt.Errorf("unknown pipeline stage %q", name)
		}
	}

	return p, nil
}
