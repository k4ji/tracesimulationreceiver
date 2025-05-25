package span

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockChildInnerCondition is a mock implementation of the Condition interface, which returns true and false one after the other.
type AtLeastConditionInnerConditionMock struct{}

func (m *AtLeastConditionInnerConditionMock) Evaluate(_ *TreeNode) (*ConditionEvaluationResult, error) {
	return NewConditionEvaluationResult([]bool{false, true, false}, true), nil
}

// TestAtLeastCondition_Evaluate_ThresholdCheck tests the AtLeastCondition with a threshold.
func TestAtLeastCondition_Evaluate(t *testing.T) {
	mock := &AtLeastConditionInnerConditionMock{}

	parent := &TreeNode{}

	condition := NewAtLeast(1, mock)

	result, err := condition.Evaluate(parent)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.evaluations, 1)
	assert.Equal(t, true, result.evaluations[0])

	condition = NewAtLeast(2, mock)

	result, err = condition.Evaluate(parent)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.evaluations, 1)
	assert.Equal(t, false, result.evaluations[0])
	assert.Equal(t, false, result.mustAggregate)
}
