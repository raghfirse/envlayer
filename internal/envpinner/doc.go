// Package envpinner provides a mechanism for pinning specific environment
// variable keys so their values cannot be overridden by subsequent merge or
// overlay operations.
//
// # Overview
//
// When building layered configurations it is sometimes necessary to lock
// certain values (e.g. security credentials, compliance-mandated settings)
// so that environment-specific overrides cannot accidentally or maliciously
// replace them.
//
// # Usage
//
//	res, err := envpinner.Pin(base, incoming, []string{"DB_PASSWORD"}, envpinner.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Vars)      // merged map with pinned values intact
//	fmt.Println(res.Violations) // keys that were silently blocked
package envpinner
