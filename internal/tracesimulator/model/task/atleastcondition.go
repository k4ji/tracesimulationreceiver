package task

// AtLeastCondition represents a condition that requires at least a certain number of nodes to meet the condition.
type AtLeastCondition struct {
	threshold int
	inner     Condition
}

// Threshold returns the minimum number of nodes that must meet the condition.
func (a *AtLeastCondition) Threshold() int {
	return a.threshold
}

// Inner returns the inner condition that must be met.
func (a *AtLeastCondition) Inner() Condition {
	return a.inner
}
