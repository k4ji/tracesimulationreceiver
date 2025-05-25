package span

import "fmt"

// ConditionEvaluationResult represents the result of evaluating a condition.
type ConditionEvaluationResult struct {
	evaluations   []bool
	mustAggregate bool
}

// NewConditionEvaluationResult creates a new ConditionEvaluationResult with the given evaluations.
func NewConditionEvaluationResult(evaluations []bool, mustAggregate bool) *ConditionEvaluationResult {
	return &ConditionEvaluationResult{
		evaluations:   evaluations,
		mustAggregate: mustAggregate,
	}
}

// Results copies and returns the evaluations of the condition evaluation.
func (r *ConditionEvaluationResult) Results() []bool {
	results := make([]bool, len(r.evaluations))
	copy(results, r.evaluations)
	return results
}

func (r *ConditionEvaluationResult) IsSatisfied() (bool, error) {
	if r.mustAggregate {
		return false, fmt.Errorf("cannot extract satisfaction from a multi-node result; wrap this with an aggregator")
	}
	return r.evaluations[0], nil
}
