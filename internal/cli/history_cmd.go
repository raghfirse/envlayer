package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nicholasgasior/envlayer/internal/envhistory"
	"github.com/nicholasgasior/envlayer/internal/loader"
)

// RunHistoryRecord loads the resolved env files and records a history entry.
func RunHistoryRecord(dir, label string, files []string, out io.Writer) error {
	if len(files) == 0 {
		return fmt.Errorf("history record: no env files specified")
	}
	if label == "" {
		return fmt.Errorf("history record: label must not be empty")
	}

	vars, err := loader.LoadFiles(files)
	if err != nil {
		return fmt.Errorf("history record: load: %w", err)
	}

	e, err := envhistory.Record(dir, label, vars)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "recorded entry %s (%s)\n", e.ID, e.Label)
	return nil
}

// RunHistoryList prints all history entries in the given directory.
func RunHistoryList(dir string, out io.Writer) error {
	entries, err := envhistory.List(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(out, "no history entries found")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tLABEL\tCREATED")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.ID, e.Label, e.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}

// RunHistoryShow prints the variables stored in a specific history entry.
func RunHistoryShow(dir, id string, out io.Writer) error {
	e, err := envhistory.Get(dir, id)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "# entry: %s  label: %s  created: %s\n",
		e.ID, e.Label, e.CreatedAt.Format("2006-01-02 15:04:05"))
	for _, k := range sortedHistoryKeys(e.Vars) {
		fmt.Fprintf(w, "%s\t=\t%s\n", k, e.Vars[k])
	}
	return w.Flush()
}

// RunHistoryDelete removes a history entry by ID.
func RunHistoryDelete(dir, id string, out io.Writer) error {
	if err := envhistory.Delete(dir, id); err != nil {
		return err
	}
	fmt.Fprintf(out, "deleted entry %s\n", id)
	return nil
}

func sortedHistoryKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	_ = strings.Join // ensure import used
	sortStringsH(keys)
	return keys
}

func sortStringsH(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// ensure os import used for potential future expansion
var _ = os.Stderr
