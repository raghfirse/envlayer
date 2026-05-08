// Package envscope provides scoped, prefix-filtered views over environment
// variable maps.
//
// A Scope lets callers work with a logical subset of environment variables
// without needing to manually filter or strip prefixes. This is useful when
// different components of an application each own a namespace of variables
// (e.g. APP_, DB_, CACHE_) and should only see their own keys.
//
// Example usage:
//
//	vars := map[string]string{
//		"DB_HOST":     "localhost",
//		"DB_PASSWORD": "secret",
//		"APP_PORT":    "8080",
//	}
//
//	scope := envscope.New("database", "DB_", vars)
//	host, _ := scope.Get("HOST")     // "localhost"
//	qualified := scope.Qualify("HOST") // "DB_HOST"
package envscope
