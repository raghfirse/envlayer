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
// # Stage Ordering
//
// Stages are executed in the order they are added via Add. The output map of
// each stage is passed as the input map to the next stage. If any stage
// returns an error, execution halts and the error is returned immediately;
// no subsequent stages are run.
//
// # Example
//
//	p := pipeline.New().
//		Add(pipeline.StageInterpolate(base)).
//		Add(pipeline.StageUppercaseKeys()).
//		Add(pipeline.StageMask(nil, ""))
//
//	result, err := p.Run(vars)
//
// # Error Handling
//
// Run returns the first error encountered along with a nil map. Callers
// should always check the returned error before using the result map.
package pipeline
