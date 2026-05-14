package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/nicholasgasior/envlayer/internal/envchain"
	"github.com/nicholasgasior/envlayer/internal/loader"
)

// RunChain loads a series of .env files in order (lowest to highest priority)
// and resolves the final merged environment, optionally filtering by key.
//
// files: ordered list of .env file paths (first = lowest priority)
// key:   if non-empty, print only that key's value
// out:   writer for output
func RunChain(files []string, key string, out io.Writer) error {
	if len(files) == 0 {
		return fmt.Errorf("chain: at least one file is required")
	}

	layers := make([]map[string]string, len(files))
	for i, f := range files {
		m, err := loader.LoadFile(f)
		if err != nil {
			return fmt.Errorf("chain: loading %q: %w", f, err)
		}
		layers[i] = m
	}

	// Reverse so that the last file provided is highest priority.
	for i, j := 0, len(layers)-1; i < j; i, j = i+1, j-1 {
		layers[i], layers[j] = layers[j], layers[i]
	}

	chain := envchain.New(layers...)

	if key != "" {
		v, err := chain.Get(key)
		if err != nil {
			return fmt.Errorf("chain: %w", err)
		}
		fmt.Fprintln(out, v)
		return nil
	}

	resolved := chain.Resolve()
	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(out, "%s=%s\n", k, resolved[k])
	}
	return nil
}

// chainFlagFiles parses a comma-separated list of file paths.
func chainFlagFiles(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// RunChainFromArgs is a thin CLI entry-point that reads flags from os.Args.
func RunChainFromArgs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envlayer chain <file1,file2,...> [key]")
	}
	files := chainFlagFiles(args[0])
	if len(files) == 0 {
		return fmt.Errorf("chain: no valid file paths provided in %q", args[0])
	}
	key := ""
	if len(args) >= 2 {
		key = args[1]
	}
	return RunChain(files, key, os.Stdout)
}
