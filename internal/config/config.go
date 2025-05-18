package config

import (
	"fmt"
	configBlueprint "github.com/k4ji/tracesimulationreceiver/internal/config/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/config/global"
)

// Config defines configuration for the Trace Simulation receiver.
type Config struct {
	// Global defines global default settings.
	Global global.Global `mapstructure:"global"`

	// Blueprint defines the blueprint of spans and their parameters.
	Blueprint configBlueprint.Blueprint `mapstructure:"blueprint"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if err := global.Validate(&cfg.Global); err != nil {
		return fmt.Errorf("global validation failed: %w", err)
	}

	if err := configBlueprint.Validate(&cfg.Blueprint); err != nil {
		return fmt.Errorf("blueprint validation failed: %w", err)
	}

	return nil
}
