// Package pipeline provides a composable transformation pipeline for
// processing environment variable maps through a series of named stages.
package pipeline

import "fmt"

// Stage is a named transformation step applied to a variable map.
type Stage struct {
	Name    string
	Apply   func(vars map[string]string) (map[string]string, error)
}

// Pipeline holds an ordered sequence of stages.
type Pipeline struct {
	stages []Stage
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a Stage to the pipeline.
func (p *Pipeline) Add(s Stage) *Pipeline {
	p.stages = append(p.stages, s)
	return p
}

// Run executes all stages in order, passing the output of each stage as
// the input to the next. The original map is never mutated.
func (p *Pipeline) Run(vars map[string]string) (map[string]string, error) {
	current := copyMap(vars)
	for _, s := range p.stages {
		result, err := s.Apply(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", s.Name, err)
		}
		current = result
	}
	return current, nil
}

// StageNames returns the names of all registered stages in order.
func (p *Pipeline) StageNames() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
