package blueprint

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/config/blueprint/service"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint"
)

const DefaultBlueprintType = "service"

// Blueprint represents the configuration for services and their span definitions.
type Blueprint struct {
	Type             string             `mapstructure:"type"`
	ServiceBlueprint *service.Blueprint `mapstructure:"service"`
}

func Validate(bp *Blueprint) error {
	switch bp.Type {
	case "service":
		if bp.ServiceBlueprint == nil {
			return fmt.Errorf("type is 'service' but service blueprint is nil")
		}
		if err := bp.ServiceBlueprint.Validate(); err != nil {
			return fmt.Errorf("service blueprint validation failed: %w", err)
		}
	}
	return nil
}

func (bp *Blueprint) To() (blueprint.Blueprint, error) {
	switch bp.Type {
	case "service":
		if bp.ServiceBlueprint == nil {
			return nil, fmt.Errorf("type is 'service' but service blueprint is nil")
		}
		return bp.ServiceBlueprint.To()
	}
	return nil, fmt.Errorf("unknown blueprint type: %s", bp.Type)
}

func Default() Blueprint {
	return Blueprint{
		Type:             DefaultBlueprintType,
		ServiceBlueprint: service.Default(),
	}
}
