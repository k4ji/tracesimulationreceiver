package span

// MarkedAsFailedCondition is a condition that checks if a node is marked as failed.
type MarkedAsFailedCondition struct{}

// NewMarkedAsFailedCondition creates a new MarkedAsFailedCondition.
func NewMarkedAsFailedCondition() MarkedAsFailedCondition {
	return MarkedAsFailedCondition{}
}

// Evaluate checks if the node is marked as failed.
func (c MarkedAsFailedCondition) Evaluate(target *TreeNode) (*ConditionEvaluationResult, error) {
	if target.Status().code == StatusCodeError {
		return NewConditionEvaluationResult([]bool{true}, false), nil
	}
	return NewConditionEvaluationResult([]bool{false}, false), nil
}
