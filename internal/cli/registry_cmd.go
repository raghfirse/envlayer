package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/envlayer/internal/loader"
	"github.com/your-org/envlayer/internal/registry"
)

// registryStore is a package-level default registry used by CLI commands.
var registryStore = registry.New()

// RunRegistryLoad reads a .env file and registers it under the given name with
// optional comma-separated tags.
//
//	RunRegistryLoad("prod", "/path/to/.env.prod", "live,stable", os.Stdout)
func RunRegistryLoad(name, filePath, tagsCSV string, out io.Writer) error {
	vars, err := loader.LoadFile(filePath)
	if err != nil {
		return fmt.Errorf("registry load: %w", err)
	}
	tags := splitCSV(tagsCSV)
	registryStore.Register(name, vars, tags...)
	fmt.Fprintf(out, "registered %q (%d keys)\n", name, len(vars))
	return nil
}

// RunRegistryGet prints all key=value pairs for the named entry.
func RunRegistryGet(name string, out io.Writer) error {
	entry, err := registryStore.Get(name)
	if err != nil {
		return err
	}
	for _, k := range sortedRegistryKeys(entry.Vars) {
		fmt.Fprintf(out, "%s=%s\n", k, entry.Vars[k])
	}
	return nil
}

// RunRegistryList prints all registered entry names, one per line.
func RunRegistryList(out io.Writer) error {
	names := registryStore.Names()
	if len(names) == 0 {
		fmt.Fprintln(out, "(no entries registered)")
		return nil
	}
	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return nil
}

// RunRegistryRemove removes the named entry from the registry.
func RunRegistryRemove(name string, out io.Writer) error {
	if err := registryStore.Remove(name); err != nil {
		return err
	}
	fmt.Fprintf(out, "removed %q\n", name)
	return nil
}

// RunRegistryFindByTag prints all entries that have the given tag.
func RunRegistryFindByTag(tag string, out io.Writer) error {
	entries := registryStore.FindByTag(tag)
	if len(entries) == 0 {
		fmt.Fprintf(out, "no entries found with tag %q\n", tag)
		return nil
	}
	for _, e := range entries {
		fmt.Fprintf(out, "%s [%s]\n", e.Name, strings.Join(e.Tags, ", "))
	}
	return nil
}

func sortedRegistryKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	_ = os.Stderr // suppress unused import
	sortStringsReg(keys)
	return keys
}

func sortStringsReg(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
