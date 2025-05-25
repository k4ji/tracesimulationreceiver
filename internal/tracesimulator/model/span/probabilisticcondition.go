package span

var _ Condition = (*Probabilistic)(nil)

// Probabilistic triggers based on a random value.
type Probabilistic struct {
	// threshold is the probability threshold for the condition to be met.
	threshold float64
	// randomness is a function that returns a random value between 0 and 1.
	randomness func() float64
}

// NewProbabilistic creates a new Probabilistic condition with the given threshold and randomness function.
func NewProbabilistic(threshold float64, randomness func() float64) *Probabilistic {
	return &Probabilistic{
		threshold:  threshold,
		randomness: randomness,
	}
}

func (p Probabilistic) Evaluate(_ *TreeNode) (*ConditionEvaluationResult, error) {
	return NewConditionEvaluationResult([]bool{p.randomness() < p.threshold}, false), nil
}
