package task

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"time"
)

type Delay struct {
	expr taskduration.Expression
}

func NewDelay(expr taskduration.Expression) (*Delay, error) {
	if expr == nil {
		panic("expr cannot be nil")
	}
	return &Delay{expr: expr}, nil
}

func (d Delay) Resolve(context interface{}) (*time.Duration, error) {
	switch d.expr.(type) {
	case *taskduration.RelativeDuration:
		parentDuration, ok := context.(*time.Duration)
		if !ok {
			return nil, fmt.Errorf("failed to resolve relative delay: invalid context type %T, expected time.Duration", context)
		}
		if parentDuration == nil {
			return nil, fmt.Errorf("relative delay requires a parent span")
		}
		delay, err := d.expr.Resolve(*parentDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve delay: %w", err)
		}
		if *delay < 0 {
			return nil, fmt.Errorf("duration cannot be negative, got %s", delay)
		}
		return delay, nil
	case *taskduration.AbsoluteDuration:
		delay, err := d.expr.Resolve(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve delay: %w", err)
		}
		if *delay < 0 {
			return nil, fmt.Errorf("duration cannot be negative, got %s", delay)
		}
		return delay, nil
	default:
		return nil, fmt.Errorf("unsupported delay type: %T", d.expr)
	}
}
