// Package validator provides validation utilities for merged environment
// variable maps produced by the envlayer pipeline.
//
// # Overview
//
// After loading and merging .env files, callers can use [Validate] to assert
// that all keys required by the application are present and non-empty.
//
// # Usage
//
//	result := validator.Validate(env, []string{"DATABASE_URL", "SECRET_KEY"})
//	if !result.IsValid() {
//		log.Fatal(result.Error())
//	}
//	for _, w := range result.Warnings {
//		log.Println("warning:", w)
//	}
package validator
