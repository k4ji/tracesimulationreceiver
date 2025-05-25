package task

// ProbabilisticCondition represents a condition that triggers based on a random value.
type ProbabilisticCondition struct {
	// threshold is the probability threshold for the condition to be met.
	threshold float64
	// randomness is a function that returns a random value between 0 and 1.
	randomness func() float64
}

// Threshold returns the probability threshold for the condition to be met.
func (p ProbabilisticCondition) Threshold() float64 {
	return p.threshold
}

// Randomness returns the randomness function used to generate random values.
func (p ProbabilisticCondition) Randomness() func() float64 {
	return p.randomness
}
