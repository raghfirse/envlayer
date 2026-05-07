// Package pipeline provides a composable, ordered processing pipeline for
// environment variable maps.
//
// A Pipeline consists of one or more named Stage values. Each Stage receives
// the output of the previous stage as its input, allowing transformations to
// be chained in a predictable, auditable order.
//
// Built-in stage constructors (StageInterpolate, StageMask, StageAddPrefix,
// StageUppercaseKeys, StageFilterPrefix) wrap existing envlayer packages so
// that common operations can be composed without boilerplate.
//
// Example:
//
//	p := pipeline.New().
//		Add(pipeline.StageInterpolate(base)).
//		Add(pipeline.StageUppercaseKeys()).
//		Add(pipeline.StageMask(nil, ""))
//
//	result, err := p.Run(vars)
package pipeline
