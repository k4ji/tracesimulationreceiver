package task

// ChildCondition represents a condition that requires a child node to meet the condition.
type ChildCondition struct {
	inner Condition
}

func (c *ChildCondition) Inner() Condition {
	return c.inner
}
