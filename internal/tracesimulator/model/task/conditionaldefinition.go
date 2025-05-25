package task

// ConditionalDefinition represents a conditional definition with a condition and an effect.
type ConditionalDefinition struct {
	condition Condition
	effects   []Effect
}

func NewConditionalDefinition(condition Condition, effects []Effect) *ConditionalDefinition {
	return &ConditionalDefinition{
		condition: condition,
		effects:   effects,
	}
}

// Condition returns the condition of the ConditionalDefinition.
func (cd *ConditionalDefinition) Condition() Condition {
	return cd.condition
}

// Effects returns the effects of the ConditionalDefinition.
func (cd *ConditionalDefinition) Effects() []Effect {
	return cd.effects
}
