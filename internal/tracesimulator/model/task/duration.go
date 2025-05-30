package task

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"time"
)

type Duration struct {
	expr taskduration.Expression
}

func NewDuration(expr taskduration.Expression) (*Duration, error) {
	if expr == nil {
		return nil, fmt.Errorf("duration expression cannot be nil")
	}
	return &Duration{expr: expr}, nil
}

func (d Duration) Resolve(context interface{}) (*time.Duration, error) {
	switch d.expr.(type) {
	case *taskduration.RelativeDuration:
		parentDuration, ok := context.(*time.Duration)
		if !ok {
			return nil, fmt.Errorf("failed to resolve relative duration: invalid context type %T, expected time.Duration", context)
		}
		if parentDuration == nil {
			return nil, fmt.Errorf("relative duration requires a parent span")
		}
		duration, err := d.expr.Resolve(*parentDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve duration: %w", err)
		}
		if *duration <= 0 {
			return nil, fmt.Errorf("duration must be greater than 0, got %s", duration)
		}
		return duration, nil
	case *taskduration.AbsoluteDuration:
		duration, err := d.expr.Resolve(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve duration: %w", err)
		}
		if *duration <= 0 {
			return nil, fmt.Errorf("duration must be greater than 0, got %s", duration)
		}
		return duration, nil
	default:
		return nil, fmt.Errorf("unsupported duration type: %T", d.expr)
	}
}
