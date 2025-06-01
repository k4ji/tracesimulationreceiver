package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
)

type Blueprint struct {
	// Default contains default values for span parameters.
	Default DefaultValues `mapstructure:"default"`
	// Services is a list of services to simulate.
	Services []Service `mapstructure:"services"`
}

// DefaultValues defines default values for span parameters.
type DefaultValues struct {
	// Delay specifies the default wait time for spans
	Delay *Delay `mapstructure:"delay"`

	// Duration specifies the default duration for spans
	Duration *Duration `mapstructure:"duration"`

	// ConditionalEffects specifies the default conditional effects for spans.
	ConditionalEffects []ConditionalEffect `mapstructure:"conditional_effects"`
}

// Validate checks the configuration for errors.
func (bp *Blueprint) Validate() error {
	refs := make(map[string]struct{})
	for _, s := range bp.Services {
		if s.Name == "" {
			return fmt.Errorf("service name cannot be empty")
		}
		for _, sd := range s.SpanDefinitions {
			if sd.Name == "" {
				return fmt.Errorf("span name cannot be empty in service %s", s.Name)
			}
			if sd.Ref != nil {
				if _, exists := refs[*sd.Ref]; exists {
					return fmt.Errorf("duplicate span ref %s found", *sd.Ref)
				} else {
					refs[*sd.Ref] = struct{}{}
				}
			}
			delay := sd.Delay.WithDefault(bp.Default.Delay)
			if delay == nil {
				return fmt.Errorf("span %s must have a delay value or global default", sd.Name)
			}
			if err := delay.ValidateAfterDefaults(); err != nil {
				return fmt.Errorf("span %s has invalid delay: %w", sd.Name, err)
			}
			duration := sd.Duration.WithDefault(bp.Default.Duration)
			if duration == nil {
				return fmt.Errorf("span %s must have a duration value or global default", sd.Name)
			}
			if err := duration.ValidateAfterDefaults(); err != nil {
				return fmt.Errorf("span %s has invalid duration: %w", sd.Name, err)
			}
		}
	}
	return nil
}

// To converts the Blueprint to a blueprint.
func (bp *Blueprint) To() (blueprint.Blueprint, error) {
	bp.prepare()
	services := make([]model.Service, len(bp.Services))
	for _, s := range bp.Services {
		s, err := s.To()
		if err != nil {
			return nil, err
		}
		services = append(services, *s)
	}
	sbp := service.NewServiceBlueprint(services)
	return &sbp, nil
}

// sets default values for span definitions if they are not specified.
func (bp *Blueprint) prepare() {
	var applyDefaults func(spanDefinitions []*SpanDefinition)
	applyDefaults = func(spanDefinitions []*SpanDefinition) {
		for _, sd := range spanDefinitions {
			sd.Delay = sd.Delay.WithDefault(bp.Default.Delay)
			sd.Duration = sd.Duration.WithDefault(bp.Default.Duration)
			childSpans := make([]*SpanDefinition, len(sd.Children))
			for i := range sd.Children {
				childSpans[i] = &sd.Children[i]
			}
			applyDefaults(childSpans)
		}
	}

	// Apply defaults to all span definitions in all services.
	for i := range bp.Services {
		sds := make([]*SpanDefinition, len(bp.Services[i].SpanDefinitions))
		for j := range bp.Services[i].SpanDefinitions {
			sds[j] = &bp.Services[i].SpanDefinitions[j]
		}
		applyDefaults(sds)
	}
}

// Default creates a new Blueprint with default values.
func Default() *Blueprint {
	return &Blueprint{
		Default: DefaultValues{
			Delay:    nil,
			Duration: nil,
		},
		Services: []Service{},
	}
}
