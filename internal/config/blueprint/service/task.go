package service

import (
	"github.com/k4ji/tracesimulator/pkg/blueprint/service/model"
	domaintask "github.com/k4ji/tracesimulator/pkg/model/task"
	"time"
)

// Task represents a task in the blueprint.
type Task struct {
	// Name is the name of the task.
	Name string `mapstructure:"name"`

	// ExternalID is an optional external identifier for the task.
	ExternalID *string `mapstructure:"id"`

	// StartAfter specifies the delay (in milliseconds) before the task starts.
	StartAfter *time.Duration `mapstructure:"startAfter"`

	// Duration specifies the duration (in milliseconds) of the task.
	Duration *time.Duration `mapstructure:"duration"`

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
	var err error
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
		StartAfter:          *t.StartAfter,
		Duration:            *t.Duration,
		Kind:                t.Kind,
		Attributes:          t.Attributes,
		Children:            children,
		ChildOf:             parentID,
		LinkedTo:            linkedTo,
		FailWithProbability: *t.FailWith.Probability,
	}, nil
}
