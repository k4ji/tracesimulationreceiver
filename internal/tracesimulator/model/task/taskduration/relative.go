package taskduration

import (
	"fmt"
	"time"
)

var _ Expression = RelativeDuration{}

// RelativeDuration represents a relative duration value, which is a percentage of a base duration.
// The value must be greater than or equal to 0.
type RelativeDuration struct {
	value float64
}

func NewRelativeDuration(value float64) (*RelativeDuration, error) {
	if value < 0 {
		return nil, fmt.Errorf("relative duration cannot be negative, got %f", value)
	}
	return &RelativeDuration{value: value}, nil
}

func (d RelativeDuration) Resolve(context interface{}) (*time.Duration, error) {
	base, ok := context.(time.Duration)
	if !ok {
		return nil, fmt.Errorf("failed to resolve relative duration: invalid context type %T, expected time.Duration", context)
	}
	if base <= 0 {
		return nil, fmt.Errorf("failed to resolve relative duration: base duration must be greater than 0, got %s", base)
	}
	r := time.Duration(float64(base.Nanoseconds()) * d.value)
	if r < 0 {
		return nil, fmt.Errorf("failed to resolve relative duration: resulting duration must be greater than or equal to 0, got %s", r)
	}
	return &r, nil
}
