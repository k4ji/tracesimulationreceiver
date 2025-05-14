package service

import (
	"fmt"
	"github.com/k4ji/tracesimulator/pkg/blueprint"
	"github.com/k4ji/tracesimulator/pkg/blueprint/service"
	"github.com/k4ji/tracesimulator/pkg/blueprint/service/model"
)

type Blueprint struct {
	// Default contains default values for task parameters.
	Default DefaultValues `mapstructure:"default"`
	// Services is a list of services to simulate.
	Services []Service `mapstructure:"services"`
}

// DefaultValues defines default values for task parameters.
type DefaultValues struct {
	// Delay specifies the default wait time for tasks
	Delay *Delay `mapstructure:"delay"`

	// Duration specifies the default duration for tasks
	Duration *Duration `mapstructure:"duration"`

	// ConditionalEffects specifies the default conditional effects for tasks.
	ConditionalEffects []ConditionalEffect `mapstructure:"conditionalEffects"`
}

// Validate checks the configuration for errors.
func (bp *Blueprint) Validate() error {
	taskIDs := make(map[string]struct{})
	for _, s := range bp.Services {
		for _, task := range s.Tasks {
			if task.ExternalID != nil {
				if _, exists := taskIDs[*task.ExternalID]; exists {
					return fmt.Errorf("duplicate task ID %s found", *task.ExternalID)
				} else {
					taskIDs[*task.ExternalID] = struct{}{}
				}
			}
			delay := task.Delay.WithDefault(bp.Default.Delay)
			if delay == nil {
				return fmt.Errorf("task %s must have a delay value or global default", task.Name)
			}
			if err := delay.ValidateAfterDefaults(); err != nil {
				return fmt.Errorf("task %s has invalid delay: %w", task.Name, err)
			}
			duration := task.Duration.WithDefault(bp.Default.Duration)
			if duration == nil {
				return fmt.Errorf("task %s must have a duration value or global default", task.Name)
			}
			if err := duration.ValidateAfterDefaults(); err != nil {
				return fmt.Errorf("task %s has invalid duration: %w", task.Name, err)
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

// sets default values for tasks if they are not specified.
func (bp *Blueprint) prepare() {
	var applyDefaults func(tasks []*Task)
	applyDefaults = func(tasks []*Task) {
		for _, task := range tasks {
			task.Delay = task.Delay.WithDefault(bp.Default.Delay)
			task.Duration = task.Duration.WithDefault(bp.Default.Duration)
			childTasks := make([]*Task, len(task.Children))
			for i := range task.Children {
				childTasks[i] = &task.Children[i]
			}
			applyDefaults(childTasks)
		}
	}

	// Apply defaults to all tasks in all services.
	for i := range bp.Services {
		tasks := make([]*Task, len(bp.Services[i].Tasks))
		for j := range bp.Services[i].Tasks {
			tasks[j] = &bp.Services[i].Tasks[j]
		}
		applyDefaults(tasks)
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
