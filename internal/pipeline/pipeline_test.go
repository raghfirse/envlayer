package pipeline_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/envlayer/internal/pipeline"
)

func TestRun_EmptyPipeline_ReturnsUnchanged(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	p := pipeline.New()
	out, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value, got %q", out["KEY"])
	}
}

func TestRun_SingleStage_Transforms(t *testing.T) {
	p := pipeline.New().Add(pipeline.Stage{
		Name: "uppercase",
		Apply: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range vars {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		},
	})
	out, err := p.Run(map[string]string{"A": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["A"])
	}
}

func TestRun_MultipleStages_ChainedCorrectly(t *testing.T) {
	addFoo := pipeline.Stage{
		Name: "add-foo",
		Apply: func(vars map[string]string) (map[string]string, error) {
			clone := make(map[string]string)
			for k, v := range vars {
				clone[k] = v
			}
			clone["FOO"] = "bar"
			return clone, nil
		},
	}
	upperVals := pipeline.Stage{
		Name: "upper-vals",
		Apply: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range vars {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		},
	}
	p := pipeline.New().Add(addFoo).Add(upperVals)
	out, err := p.Run(map[string]string{"X": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "BAR" {
		t.Errorf("expected BAR, got %q", out["FOO"])
	}
	if out["X"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["X"])
	}
}

func TestRun_StageError_WrapsWithStageName(t *testing.T) {
	p := pipeline.New().Add(pipeline.Stage{
		Name: "fail-stage",
		Apply: func(vars map[string]string) (map[string]string, error) {
			return nil, errors.New("something went wrong")
		},
	})
	_, err := p.Run(map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fail-stage") {
		t.Errorf("expected stage name in error, got: %v", err)
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "original"}
	p := pipeline.New().Add(pipeline.Stage{
		Name: "mutate",
		Apply: func(vars map[string]string) (map[string]string, error) {
			vars["KEY"] = "mutated"
			return vars, nil
		},
	})
	p.Run(input) //nolint
	if input["KEY"] != "original" {
		t.Errorf("input was mutated: got %q", input["KEY"])
	}
}

func TestStageNames_ReturnsInOrder(t *testing.T) {
	p := pipeline.New().
		Add(pipeline.Stage{Name: "first", Apply: func(v map[string]string) (map[string]string, error) { return v, nil }}).
		Add(pipeline.Stage{Name: "second", Apply: func(v map[string]string) (map[string]string, error) { return v, nil }})
	names := p.StageNames()
	if len(names) != 2 || names[0] != "first" || names[1] != "second" {
		t.Errorf("unexpected stage names: %v", names)
	}
}
