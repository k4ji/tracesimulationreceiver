package span

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Condition is an interface for evaluating whether a condition is met.
type Condition interface {
	// Evaluate evaluates the condition based on the provided context and returns true if the condition is met.
	Evaluate(target *TreeNode) (*ConditionEvaluationResult, error)
}

// FromConditionSpec converts a Condition spec to a Condition.
func FromConditionSpec(spec task.Condition) (Condition, error) {
	switch spec.Kind() {
	case task.ConditionKindProbabilistic:
		if spec.Probabilistic() == nil {
			return nil, fmt.Errorf("probabilistic condition requires a probability")
		}
		return NewProbabilistic(
			spec.Probabilistic().Threshold(),
			spec.Probabilistic().Randomness(),
		), nil
	case task.ConditionKindHasAttribute:
		if spec.HasAttribute() == nil {
			return nil, fmt.Errorf("hasAttribute condition requires a key")
		}
		return NewHasAttribute(spec.HasAttribute().Key()), nil
	case task.ConditionKindChild:
		if spec.Child() == nil {
			return nil, fmt.Errorf("child condition requires a child condition")
		}
		innerCondition, err := FromConditionSpec(spec.Child().Inner())
		if err != nil {
			return nil, fmt.Errorf("failed to convert inner condition: %w", err)
		}
		return NewChild(innerCondition), nil
	case task.ConditionKindAtLeast:
		if spec.AtLeast() == nil {
			return nil, fmt.Errorf("atLeast condition requires a count")
		}
		innerCondition, err := FromConditionSpec(spec.AtLeast().Inner())
		if err != nil {
			return nil, fmt.Errorf("failed to convert inner condition: %w", err)
		}
		return NewAtLeast(
			spec.AtLeast().Threshold(),
			innerCondition,
		), nil
	case task.ConditionKindMarkedAsFailed:
		if spec.MarkedAsFailed() == nil {
			return nil, fmt.Errorf("markedAsFailed condition requires a message")
		}
		return NewMarkedAsFailedCondition(), nil
	default:
		return nil, fmt.Errorf("unsupported condition type: %s", spec.Kind())
	}
}
