package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envlayer/internal/differ"
	"github.com/user/envlayer/internal/loader"
)

// DiffOptions configures the RunDiff command.
type DiffOptions struct {
	FromFile string
	ToFile   string
	FromLabel string
	ToLabel   string
	Out      io.Writer
}

// RunDiff loads two .env files and prints a human-readable diff to Out.
func RunDiff(opts DiffOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.FromLabel == "" {
		opts.FromLabel = opts.FromFile
	}
	if opts.ToLabel == "" {
		opts.ToLabel = opts.ToFile
	}

	from, err := loader.LoadFile(opts.FromFile)
	if err != nil {
		return fmt.Errorf("loading from-file %q: %w", opts.FromFile, err)
	}

	to, err := loader.LoadFile(opts.ToFile)
	if err != nil {
		return fmt.Errorf("loading to-file %q: %w", opts.ToFile, err)
	}

	result := differ.Diff(from, to, opts.FromLabel, opts.ToLabel)

	fmt.Fprintln(opts.Out, differ.Summary(result))

	for _, c := range result.Changes {
		switch c.Kind {
		case differ.Added:
			fmt.Fprintf(opts.Out, "  + %s=%q\n", c.Key, c.NewValue)
		case differ.Removed:
			fmt.Fprintf(opts.Out, "  - %s=%q\n", c.Key, c.OldValue)
		case differ.Changed:
			fmt.Fprintf(opts.Out, "  ~ %s: %q → %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}

	if len(result.Changes) == 0 {
		fmt.Fprintln(opts.Out, "  (no changes)")
	}

	return nil
}
