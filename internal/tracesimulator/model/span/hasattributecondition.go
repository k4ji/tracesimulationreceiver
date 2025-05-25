package span

// HasAttribute is a condition that checks if a node has a specific attribute.
type HasAttribute struct {
	key string
}

func NewHasAttribute(key string) HasAttribute {
	return HasAttribute{
		key: key,
	}
}

func (c HasAttribute) Evaluate(target *TreeNode) (*ConditionEvaluationResult, error) {
	_, ok := target.Attributes()[c.key]
	return NewConditionEvaluationResult([]bool{ok}, false), nil
}
