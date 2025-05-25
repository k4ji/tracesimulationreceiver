package task

type ConditionKind string

const (
	ConditionKindProbabilistic  ConditionKind = "probabilistic"
	ConditionKindAtLeast        ConditionKind = "atLeast"
	ConditionKindChild          ConditionKind = "child"
	ConditionKindHasAttribute   ConditionKind = "hasAttribute"
	ConditionKindMarkedAsFailed ConditionKind = "markedAsFailed"
)

// Condition is an interface for evaluating whether an effect should be applied.
type Condition struct {
	// kind is the type of condition.
	kind ConditionKind
	// probabilistic is the probability of the condition being met.
	probabilistic *ProbabilisticCondition
	// atLeast is the minimum number of nodes that must meet the condition.
	atLeast *AtLeastCondition
	// child is the child condition that must be met.
	child *ChildCondition
	// hasAttribute is the attribute that must be present.
	hasAttribute *HasAttributeCondition
	// markedAsFailed is the condition that checks if the task is marked as failed.
	markedAsFailed *MarkedAsFailedCondition
}

// NewProbabilisticCondition creates a new Condition with the given probability.
func NewProbabilisticCondition(threshold float64, randomness func() float64) Condition {
	return Condition{
		kind: ConditionKindProbabilistic,
		probabilistic: &ProbabilisticCondition{
			threshold:  threshold,
			randomness: randomness,
		},
	}
}

// NewAtLeastCondition creates a new Condition with the given threshold and inner condition.
func NewAtLeastCondition(threshold int, inner Condition) Condition {
	return Condition{
		kind: ConditionKindAtLeast,
		atLeast: &AtLeastCondition{
			threshold: threshold,
			inner:     inner,
		},
	}
}

// NewChildCondition creates a new Condition with the given inner condition.
func NewChildCondition(inner Condition) Condition {
	return Condition{
		kind: ConditionKindChild,
		child: &ChildCondition{
			inner: inner,
		},
	}
}

// NewHasAttributeCondition creates a new Condition with the given status code.
func NewHasAttributeCondition(key string) Condition {
	return Condition{
		kind: ConditionKindHasAttribute,
		hasAttribute: &HasAttributeCondition{
			key: key,
		},
	}
}

// NewMarkedAsFailedCondition creates a new Condition that checks if the task is marked as failed.
func NewMarkedAsFailedCondition() Condition {
	return Condition{
		kind:           ConditionKindMarkedAsFailed,
		markedAsFailed: &MarkedAsFailedCondition{},
	}
}

func (c Condition) Kind() ConditionKind {
	return c.kind
}

func (c Condition) Probabilistic() *ProbabilisticCondition {
	return c.probabilistic
}

func (c Condition) AtLeast() *AtLeastCondition {
	return c.atLeast
}

func (c Condition) Child() *ChildCondition {
	return c.child
}

func (c Condition) HasAttribute() *HasAttributeCondition {
	return c.hasAttribute
}

func (c Condition) MarkedAsFailed() *MarkedAsFailedCondition {
	return c.markedAsFailed
}
