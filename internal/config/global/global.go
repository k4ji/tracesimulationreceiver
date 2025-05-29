package global

import (
	"fmt"
	"time"
)

const DefaultInterval = 5 * time.Second
const DefaultEndTimeOffset = 0 * time.Second

// Global defines global default settings for span definitions and intervals.
type Global struct {
	// Interval specifies the time interval at which traces are generated (e.g., "5s").
	Interval time.Duration `mapstructure:"interval"`
	// EndTimeOffset specifies the base offset for the end time of spans.
	EndTimeOffset time.Duration `mapstructure:"end_time_offset"`
}

func Validate(g *Global) error {
	if g.Interval <= 0 {
		return fmt.Errorf("global interval must be greater than 0")
	}
	// EndTimeOffset must be between 1 year in the past and 1 year in the future
	if g.EndTimeOffset < -365*24*time.Hour || g.EndTimeOffset > 365*24*time.Hour {
		return fmt.Errorf("global end_time_offset must be between -1 year and +1 year")
	}
	return nil
}

func Default() Global {
	return Global{
		Interval:      DefaultInterval,
		EndTimeOffset: DefaultEndTimeOffset,
	}
}
