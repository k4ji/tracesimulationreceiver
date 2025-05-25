package span

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockChildInnerCondition is a mock implementation of the Condition interface, which returns true and false one after the other.
type ChildConditionInnerConditionMock struct {
	evaluationCount int
}

func (m *ChildConditionInnerConditionMock) Evaluate(target *TreeNode) (*ConditionEvaluationResult, error) {
	m.evaluationCount++
	if m.evaluationCount%2 == 0 {
		return NewConditionEvaluationResult([]bool{true}, false), nil
	}
	return NewConditionEvaluationResult([]bool{false}, false), nil
}

func TestChildCondition_Evaluate(t *testing.T) {
	mock := &ChildConditionInnerConditionMock{evaluationCount: 0}

	child1 := &TreeNode{}
	child2 := &TreeNode{}
	child3 := &TreeNode{}
	parent := &TreeNode{
		children: []*TreeNode{child1, child2, child3},
	}

	condition := NewChild(mock)

	result, err := condition.Evaluate(parent)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.evaluations, 3)
	assert.Equal(t, false, result.evaluations[0])
	assert.Equal(t, true, result.evaluations[1])
	assert.Equal(t, false, result.evaluations[2])
	assert.Equal(t, true, result.mustAggregate)
}

func TestChildCondition_Evaluate_Nested(t *testing.T) {
	mock := &ChildConditionInnerConditionMock{evaluationCount: 0}

	node := &TreeNode{
		children: []*TreeNode{
			{},
			{
				children: []*TreeNode{
					{},
					{},
				},
			},
			{},
		},
	}

	condition := NewChild(NewChild(mock))

	result, err := condition.Evaluate(node)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.evaluations, 2)
	assert.Equal(t, false, result.evaluations[0])
	assert.Equal(t, true, result.evaluations[1])
	assert.Equal(t, true, result.mustAggregate)
}
