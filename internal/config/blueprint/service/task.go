package service

import (
	"github.com/k4ji/tracesimulator/pkg/blueprint/service/model"
	domaintask "github.com/k4ji/tracesimulator/pkg/model/task"
)

// Task represents a task in the blueprint.
type Task struct {
	// Name is the name of the task.
	Name string `mapstructure:"name"`

	// ExternalID is an optional external identifier for the task.
	ExternalID *string `mapstructure:"id"`

	// Delay specifies the delay in duration or relative duration to parent duration before the task starts.
	Delay *Delay `mapstructure:"delay"`

	// Duration specifies the duration of the task.
	Duration *Duration `mapstructure:"duration"`

	// Kind specifies the type or category of the task (e.g., "client", "server").
	Kind string `mapstructure:"kind"`

	// Attributes contains optional attributes for the task.
	Attributes map[string]string `mapstructure:"attributes"`

	// Children is a list of child tasks triggered by this task.
	Children []Task `mapstructure:"children"`

	// ChildOf is an optional parent task identifier.
	ChildOf *string `mapstructure:"childOf"`

	// LinkedTo is a list of task identifiers this task is linked to.
	LinkedTo []*string `mapstructure:"linkedTo"`

	// FailWith specifies failure conditions for tasks.
	FailWith FailureCondition `mapstructure:"failWith"`
}

// To return model.Task
func (t *Task) To() (*model.Task, error) {
	var externalID *domaintask.ExternalID
	var parentID *domaintask.ExternalID
	var linkedTo []*domaintask.ExternalID
	var children []model.Task
	delay, err := t.Delay.To()
	if err != nil {
		return nil, err
	}
	duration, err := t.Duration.To()
	if err != nil {
		return nil, err
	}
	if t.ExternalID != nil {
		externalID, err = domaintask.NewExternalID(*t.ExternalID)
		if err != nil {
			return nil, err
		}
	}
	if t.ChildOf != nil {
		parentID, err = domaintask.NewExternalID(*t.ChildOf)
		if err != nil {
			return nil, err
		}
	}
	if t.LinkedTo != nil {
		linkedTo = make([]*domaintask.ExternalID, len(t.LinkedTo))
		for j, link := range t.LinkedTo {
			linkedTo[j], err = domaintask.NewExternalID(*link)
			if err != nil {
				return nil, err
			}
		}
	}
	if t.Children != nil {
		children = make([]model.Task, len(t.Children))
		for i, child := range t.Children {
			c, err := child.To()
			if err != nil {
				return nil, err
			}
			children[i] = *c
		}
	}
	return &model.Task{
		Name:                t.Name,
		ExternalID:          externalID,
		Delay:               *delay,
		Duration:            *duration,
		Kind:                t.Kind,
		Attributes:          t.Attributes,
		Children:            children,
		ChildOf:             parentID,
		LinkedTo:            linkedTo,
		FailWithProbability: *t.FailWith.Probability,
	}, nil
}
