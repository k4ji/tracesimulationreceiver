package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/config/utils"
	"github.com/k4ji/tracesimulator/pkg/model/task"
	"strconv"
	"time"
)

// Duration represents a duration configuration in the blueprint
type Duration struct {
	// Value is the duration value
	Value *string `mapstructure:"for"`
	// Mode is the duration mode (absolute or relative)
	Mode *string `mapstructure:"as"`
}

// To converts the duration to a model.Value
func (d *Duration) To() (*task.Duration, error) {
	if err := d.ValidateAfterDefaults(); err != nil {
		return nil, err
	}
	td := TaskDuration{
		Duration: *d.Value,
		Mode:     FromString(*d.Mode),
	}
	expr, err := td.To()
	if err != nil {
		return nil, fmt.Errorf("failed to convert duration: %w", err)
	}
	duration, err := task.NewDuration(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to create duration: %w", err)
	}
	return duration, nil
}

// ValidateAfterDefaults checks if the duration is valid.
// This is intentionally not named Validate to avoid automatic calls by `xconfmap`
// before default values are applied. If needed, consider introducing a separate type
// for pre-default unmarshaling.
func (d *Duration) ValidateAfterDefaults() error {
	if d == nil {
		return fmt.Errorf("missing duration")
	}
	if d.Value == nil {
		return fmt.Errorf("missing required field: duration.for")
	}
	if d.Mode == nil {
		return fmt.Errorf("missing required field: duration.mode")
	}
	switch FromString(*d.Mode) {
	case AbsoluteMode:
		dur, err := time.ParseDuration(*d.Value)
		if err != nil {
			return fmt.Errorf("invalid absolute duration format: %w", err)
		}
		if dur <= 0 {
			return fmt.Errorf("absolute duration must be greater than 0")
		}
	case RelativeMode:
		f, err := strconv.ParseFloat(*d.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid relative duration format: %w", err)
		}
		if f <= 0 {
			return fmt.Errorf("relative duration must be greater than 0")
		}
	default:
		return fmt.Errorf("unsupported duration mode: %s", *d.Mode)
	}
	return nil
}

// WithDefault returns a new Duration with default values applied.
func (d *Duration) WithDefault(dd *Duration) *Duration {
	if d == nil {
		return dd
	}
	if dd == nil {
		return d
	}
	return &Duration{
		Value: utils.Coalesce(d.Value, dd.Value),
		Mode:  utils.Coalesce(d.Mode, dd.Mode),
	}
}
