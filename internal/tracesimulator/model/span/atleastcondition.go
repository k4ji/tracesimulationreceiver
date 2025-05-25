package span

var _ Condition = (*AtLeastCondition)(nil)

// AtLeastCondition represents a condition that requires at least a certain number of nodes to meet the condition.
type AtLeastCondition struct {
	// threshold is the minimum number of nodes that must meet the condition.
	threshold int
	// inner is the inner condition that must be met.
	inner Condition
}

func NewAtLeast(threshold int, inner Condition) Condition {
	return AtLeastCondition{threshold: threshold, inner: inner}
}

func (c AtLeastCondition) Evaluate(target *TreeNode) (*ConditionEvaluationResult, error) {
	count := 0
	cr, err := c.inner.Evaluate(target)
	if err != nil {
		return nil, err
	}
	rs := cr.Results()
	for i := 0; i < len(rs); i++ {
		if rs[i] {
			count++
			if count >= c.threshold {
				return NewConditionEvaluationResult([]bool{true}, false), nil
			}
		}
	}
	return NewConditionEvaluationResult([]bool{false}, false), nil
}
