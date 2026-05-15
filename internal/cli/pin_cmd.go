package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/your-org/envlayer/internal/envpinner"
	"github.com/your-org/envlayer/internal/loader"
)

// RunPin loads one or more .env files and re-emits them with the supplied keys
// pinned — i.e. values from the first (base) file win for those keys.
//
// flags:
//
//	--files   comma-separated list of .env files (first = base)
//	--pin     comma-separated list of keys to pin
//	--strict  return an error if a pinned key is overridden with a different value
func RunPin(args []string, out io.Writer) error {
	files, pin, strict, err := parsePinArgs(args)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("pin: at least one --files entry is required")
	}
	if len(pin) == 0 {
		return fmt.Errorf("pin: at least one --pin key is required")
	}

	base, err := loader.LoadFile(files[0])
	if err != nil {
		return fmt.Errorf("pin: loading base file: %w", err)
	}

	opts := envpinner.Options{StrictMode: strict}
	current := base

	for _, f := range files[1:] {
		incoming, err := loader.LoadFile(f)
		if err != nil {
			return fmt.Errorf("pin: loading file %q: %w", f, err)
		}
		res, err := envpinner.Pin(current, incoming, pin, opts)
		if err != nil {
			return fmt.Errorf("pin: %w", err)
		}
		if len(res.Violations) > 0 {
			fmt.Fprintf(os.Stderr, "pin: blocked override for pinned keys: %s\n",
				strings.Join(res.Violations, ", "))
		}
		current = res.Vars
	}

	keys := make([]string, 0, len(current))
	for k := range current {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(out, "%s=%s\n", k, current[k])
	}
	return nil
}

func parsePinArgs(args []string) (files, pin []string, strict bool, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--files":
			i++
			if i >= len(args) {
				return nil, nil, false, fmt.Errorf("pin: --files requires a value")
			}
			files = splitCSV(args[i])
		case "--pin":
			i++
			if i >= len(args) {
				return nil, nil, false, fmt.Errorf("pin: --pin requires a value")
			}
			pin = splitCSV(args[i])
		case "--strict":
			strict = true
		}
	}
	return
}
