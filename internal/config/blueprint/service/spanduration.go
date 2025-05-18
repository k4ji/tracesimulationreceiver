package service

import (
	"fmt"
	"github.com/k4ji/tracesimulator/pkg/model/task/taskduration"
	"strconv"
	"time"
)

// Mode represents the mode of span duration
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

// SpanDuration represents a span duration configuration
type SpanDuration struct {
	// Duration is the duration string
	Duration string
	// Mode is the duration mode (absolute or relative)
	Mode Mode
}

// To converts the SpanDuration to a taskduration.Expression
func (d *SpanDuration) To() (taskduration.Expression, error) {
	switch d.Mode {
	case AbsoluteMode:
		dur, err := time.ParseDuration(d.Duration)
		if err != nil {
			return nil, fmt.Errorf("invalid absolute span duration: %w", err)
		}
		expr, err := taskduration.NewAbsoluteDuration(dur)
		if err != nil {
			return nil, fmt.Errorf("failed to create absolute span duration: %w", err)
		}
		return expr, nil
	case RelativeMode:
		f, err := strconv.ParseFloat(d.Duration, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid relative span duration: %w", err)
		}
		expr, err := taskduration.NewRelativeDuration(f)
		if err != nil {
			return nil, fmt.Errorf("failed to create relative span duration: %w", err)
		}
		return expr, nil
	case UnknownMode:
		return nil, fmt.Errorf("unsupported span duration mode: %s", d.Mode)
	}
	return nil, fmt.Errorf("invalid span duration mode: %s", d.Mode)
}
