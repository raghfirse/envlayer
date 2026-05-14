package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourorg/envlayer/internal/envaliaser"
	"github.com/yourorg/envlayer/internal/loader"
)

// AliasArgs holds configuration for the alias command.
type AliasArgs struct {
	// Files is the ordered list of .env files to load.
	Files []string
	// Aliases is a list of "SRC=ALIAS1,ALIAS2" strings.
	Aliases []string
	// KeepOriginal retains the source key in the output.
	KeepOriginal bool
	// FailOnMissing returns an error if a source key is absent.
	FailOnMissing bool
	// Out is the writer for the rendered output (defaults to os.Stdout).
	Out io.Writer
}

// RunAlias loads env files, applies the provided alias mappings, and prints
// the resulting key=value pairs as a dotenv-style output.
func RunAlias(args AliasArgs) error {
	if len(args.Files) == 0 {
		return fmt.Errorf("alias: at least one env file is required")
	}
	if len(args.Aliases) == 0 {
		return fmt.Errorf("alias: at least one alias mapping is required")
	}

	out := args.Out
	if out == nil {
		out = os.Stdout
	}

	vars, err := loader.LoadFiles(args.Files)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	aliasMap, err := parseAliasArgs(args.Aliases)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	opts := envaliaser.Options{
		KeepOriginal:  args.KeepOriginal,
		FailOnMissing: args.FailOnMissing,
	}
	result, err := envaliaser.Apply(vars, aliasMap, opts)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

// parseAliasArgs converts ["SRC=A,B"] into an AliasMap.
func parseAliasArgs(raw []string) (envaliaser.AliasMap, error) {
	am := make(envaliaser.AliasMap, len(raw))
	for _, entry := range raw {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid alias mapping %q: expected SRC=ALIAS1[,ALIAS2,...]", entry)
		}
		src := strings.TrimSpace(parts[0])
		targets := splitCSV(parts[1])
		if len(targets) == 0 {
			return nil, fmt.Errorf("invalid alias mapping %q: no targets specified", entry)
		}
		am[src] = targets
	}
	return am, nil
}
