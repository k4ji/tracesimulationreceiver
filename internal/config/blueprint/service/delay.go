package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/config/utils"
	"github.com/k4ji/tracesimulator/pkg/model/task"
	"strconv"
	"time"
)

// Delay represents a delay configuration in the blueprint
type Delay struct {
	// Value is the delay value
	Value *string `mapstructure:"for"`
	// Mode is the delay mode (absolute or relative)
	Mode *string `mapstructure:"as"`
}

// To converts the delay to a model.Delay
func (d *Delay) To() (*task.Delay, error) {
	if err := d.ValidateAfterDefaults(); err != nil {
		return nil, err
	}
	td := SpanDuration{
		Duration: *d.Value,
		Mode:     FromString(*d.Mode),
	}
	expr, err := td.To()
	if err != nil {
		return nil, fmt.Errorf("failed to convert delay: %w", err)
	}
	delay, err := task.NewDelay(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to create delay: %w", err)
	}
	return delay, nil
}

// ValidateAfterDefaults checks if the delay is valid.
// This is intentionally not named Validate to avoid automatic calls by `xconfmap`
// before default values are applied. If needed, consider introducing a separate type
// for pre-default unmarshaling.
func (d *Delay) ValidateAfterDefaults() error {
	if d == nil {
		return fmt.Errorf("missing delay")
	}
	if d.Value == nil {
		return fmt.Errorf("missing required field: delay.for")
	}
	if d.Mode == nil {
		return fmt.Errorf("missing required field: delay.mode")
	}
	switch FromString(*d.Mode) {
	case AbsoluteMode:
		dur, err := time.ParseDuration(*d.Value)
		if err != nil {
			return fmt.Errorf("invalid absolute delay format: %w", err)
		}
		if dur < 0 {
			return fmt.Errorf("absolute delay must be non-negative")
		}
	case RelativeMode:
		f, err := strconv.ParseFloat(*d.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid relative delay format: %w", err)
		}
		if f < 0 {
			return fmt.Errorf("relative delay must be non-negative")
		}
	default:
		return fmt.Errorf("unsupported delay mode: %s", *d.Mode)
	}
	return nil
}

// WithDefault returns a new Delay with default values applied.
func (d *Delay) WithDefault(dd *Delay) *Delay {
	if d == nil {
		return dd
	}
	if dd == nil {
		return d
	}
	return &Delay{
		Value: utils.Coalesce(d.Value, dd.Value),
		Mode:  utils.Coalesce(d.Mode, dd.Mode),
	}
}
