package global

import (
	"fmt"
	"time"
)

const DefaultInterval = 5 * time.Second

// Global defines global default settings for span definitions and intervals.
type Global struct {
	// Interval specifies the time interval at which traces are generated (e.g., "5s").
	Interval time.Duration `mapstructure:"interval"`
}

func Validate(g *Global) error {
	if g.Interval <= 0 {
		return fmt.Errorf("global interval must be greater than 0")
	}
	return nil
}

func Default() Global {
	return Global{
		Interval: DefaultInterval,
	}
}
