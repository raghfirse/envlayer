// Package schema defines a JSON-based schema format for declaring expected
// environment variables, their types, required status, defaults, and
// human-readable descriptions.
//
// A schema file is a JSON object where each key is a variable name and the
// value is a Field descriptor:
//
//	{
//	  "PORT":  { "type": "int",    "required": true,  "description": "HTTP port" },
//	  "DEBUG": { "type": "bool",   "required": false, "default": "false" },
//	  "NAME":  { "type": "string", "required": true }
//	}
//
// Supported types: string, int, bool, float.
//
// Use LoadSchema to read a schema from disk, and Validate to check a
// map[string]string of resolved environment variables against it.
package schema
