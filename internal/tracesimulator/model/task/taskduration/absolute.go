package taskduration

import (
	"fmt"
	"time"
)

var _ Expression = AbsoluteDuration{}

// AbsoluteDuration represents an absolute duration.
// duration cannot be negative.
type AbsoluteDuration struct {
	duration time.Duration
}

func NewAbsoluteDuration(duration time.Duration) (*AbsoluteDuration, error) {
	if duration < 0 {
		return nil, fmt.Errorf("absolute duration cannot be negative, got %s", duration)
	}
	return &AbsoluteDuration{duration: duration}, nil
}

func (f AbsoluteDuration) Resolve(_ interface{}) (*time.Duration, error) {
	return &f.duration, nil
}
