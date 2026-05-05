// Package printer provides utilities for rendering merged environment
// variable maps to an output destination in a specified format.
//
// It acts as the final stage in the envlayer pipeline, sitting on top of
// the exporter package and adding destination management and key filtering.
//
// Supported formats mirror those of the exporter package:
//   - "shell"  — export KEY=VALUE statements suitable for shell sourcing
//   - "dotenv" — KEY=VALUE pairs in .env file format
//   - "json"   — a JSON object mapping keys to values
//
// Example usage:
//
//	p := printer.New(printer.Options{
//		Format: "dotenv",
//		Prefix: "APP_",
//	})
//	p.Print(mergedEnv)
package printer
