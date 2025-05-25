package task

// HasAttributeCondition is a condition that checks if a task has a specific attribute.
type HasAttributeCondition struct {
	key string
}

func (c HasAttributeCondition) Key() string {
	return c.key
}
