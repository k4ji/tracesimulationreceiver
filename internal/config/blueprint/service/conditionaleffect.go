package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// ConditionalEffect represents a conditional effect that can be applied to a span.
type ConditionalEffect struct {
	// Condition is the condition that must be met for the effect to be applied.
	Condition Condition `mapstructure:"condition"`
	// Effect is the effect to be applied if the condition is met.
	Effects []Effect `mapstructure:"effects"`
}

// To creates a new ConditionalEffect with the given kind and attributes.
func (c *ConditionalEffect) To() (*task.ConditionalDefinition, error) {
	condition, err := c.Condition.To()
	if err != nil {
		return nil, fmt.Errorf("failed to convert condition: %w", err)
	}

	var effects []task.Effect
	for _, effect := range c.Effects {
		e, err := effect.To()
		if err != nil {
			return nil, fmt.Errorf("failed to convert effect: %w", err)
		}
		effects = append(effects, *e)
	}

	def := task.NewConditionalDefinition(
		*condition,
		effects,
	)
	return def, nil
}
