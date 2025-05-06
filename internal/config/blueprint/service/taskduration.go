package service

import (
	"fmt"
	"github.com/k4ji/tracesimulator/pkg/model/task/taskduration"
	"strconv"
	"time"
)

// Mode represents the mode of task duration
type Mode string

const (
	// AbsoluteMode represents an absolute duration
	AbsoluteMode Mode = "absolute"
	// RelativeMode represents a relative duration
	RelativeMode Mode = "relative"
	// UnknownMode represents an unknown duration
	UnknownMode Mode = "unknown"
)

func FromString(mode string) Mode {
	switch mode {
	case "absolute":
		return AbsoluteMode
	case "relative":
		return RelativeMode
	default:
		return UnknownMode
	}
}

// TaskDuration represents a task duration configuration
type TaskDuration struct {
	// Duration is the duration string
	Duration string
	// Mode is the duration mode (absolute or relative)
	Mode Mode
}

// To converts the TaskDuration to an Expression
func (d *TaskDuration) To() (taskduration.Expression, error) {
	switch d.Mode {
	case AbsoluteMode:
		dur, err := time.ParseDuration(d.Duration)
		if err != nil {
			return nil, fmt.Errorf("invalid absolute task duration: %w", err)
		}
		expr, err := taskduration.NewAbsoluteDuration(dur)
		if err != nil {
			return nil, fmt.Errorf("failed to create absolute task duration: %w", err)
		}
		return expr, nil
	case RelativeMode:
		f, err := strconv.ParseFloat(d.Duration, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid relative task duration: %w", err)
		}
		expr, err := taskduration.NewRelativeDuration(f)
		if err != nil {
			return nil, fmt.Errorf("failed to create relative task duration: %w", err)
		}
		return expr, nil
	case UnknownMode:
		return nil, fmt.Errorf("unsupported task duration mode: %s", d.Mode)
	}
	return nil, fmt.Errorf("invalid task duration mode: %s", d.Mode)
}
