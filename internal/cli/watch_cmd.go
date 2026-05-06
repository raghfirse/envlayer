package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/yourorg/envlayer/internal/loader"
	"github.com/yourorg/envlayer/internal/masker"
	"github.com/yourorg/envlayer/internal/resolver"
	"github.com/yourorg/envlayer/internal/watcher"
)

// WatchOptions configures the watch sub-command.
type WatchOptions struct {
	Dir       string
	Env       string
	Interval  time.Duration
	MaskKeys  []string
	Quiet     bool
}

// RunWatch resolves env files, prints the initial merged state, then
// watches for changes and re-prints whenever a file is modified.
func RunWatch(opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = 500 * time.Millisecond
	}

	paths, err := resolver.Resolve(opts.Dir, opts.Env)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	printMerged := func() error {
		vars, err := loader.LoadFiles(paths)
		if err != nil {
			return err
		}
		if len(opts.MaskKeys) > 0 {
			vars = masker.Mask(vars, masker.Options{SensitiveKeys: opts.MaskKeys})
		}
		for k, v := range vars {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	if !opts.Quiet {
		fmt.Fprintln(os.Stderr, "[envlayer] watching:", paths)
	}

	if err := printMerged(); err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	events := watcher.Watch(paths, opts.Interval, done)

	for ev := range events {
		if !opts.Quiet {
			fmt.Fprintf(os.Stderr, "[envlayer] %s: %s\n", ev.Kind, ev.Path)
		}
		if err := printMerged(); err != nil {
			fmt.Fprintln(os.Stderr, "[envlayer] error:", err)
		}
	}

	return nil
}
