package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"math/rand"
	"time"
)

// Condition represents a condition that determines whether an effect should be applied.
type Condition struct {
	// Kind is the type of condition.
	Kind string `mapstructure:"kind"`
	// Probabilistic is the probability of the condition being met.
	Probabilistic *Probabilistic `mapstructure:"probabilistic"`
}

// Probabilistic represents a probabilistic condition.
type Probabilistic struct {
	// Threshold is the threshold for the condition.
	Threshold float64 `mapstructure:"threshold"`
}

// To converts the condition to a domain model.
func (c *Condition) To() (*task.Condition, error) {
	switch c.Kind {
	case "probabilistic":
		if c.Probabilistic == nil {
			return nil, fmt.Errorf("probabilistic condition requires a threshold")
		}
		if c.Probabilistic.Threshold < 0 || c.Probabilistic.Threshold > 1 {
			return nil, fmt.Errorf("probabilistic condition threshold must be between 0 and 1")
		}
		condition := task.NewProbabilisticCondition(c.Probabilistic.Threshold, rand.New(rand.NewSource(time.Now().UnixNano())).Float64)
		return &condition, nil
	case "childMarkedAsFailed":
		condition := task.NewAtLeastCondition(1, task.NewChildCondition(task.NewMarkedAsFailedCondition()))
		return &condition, nil
	default:
		return nil, fmt.Errorf("unknown condition type: %s", c.Kind)
	}
}
