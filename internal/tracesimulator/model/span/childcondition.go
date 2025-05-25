package span

var _ Condition = (*ChildCondition)(nil)

// ChildCondition represents a condition that requires child nodes to meet the condition.
type ChildCondition struct {
	inner Condition
}

func NewChild(inner Condition) Condition {
	return ChildCondition{inner: inner}
}

func (c ChildCondition) Evaluate(target *TreeNode) (*ConditionEvaluationResult, error) {
	var results []bool
	for _, child := range target.Children() {
		cr, err := c.inner.Evaluate(child)
		if err != nil {
			return nil, err
		}
		rs := cr.Results()
		results = append(results, rs...)
	}
	return NewConditionEvaluationResult(results, true), nil
}
