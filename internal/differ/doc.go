// Package differ computes structured diffs between two environment variable
// maps, labelled by their source (e.g. "dev" and "prod").
//
// Basic usage:
//
//	from := map[string]string{"HOST": "localhost", "PORT": "5432"}
//	to   := map[string]string{"HOST": "db.prod",   "PORT": "5432", "TLS": "true"}
//
//	result := differ.Diff(from, to, "dev", "prod")
//	fmt.Println(differ.Summary(result))
//	// dev → prod: +1 -0 ~1
//
//	for _, c := range result.Changes {
//		fmt.Printf("%s [%s] %q → %q\n", c.Key, c.Kind, c.OldValue, c.NewValue)
//	}
//
// Changes are sorted by key for deterministic output.
package differ
