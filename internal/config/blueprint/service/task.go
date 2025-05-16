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

	// Events is a list of events associated with the task.
	Events []Event `mapstructure:"events"`

	// Parent is an optional parent task identifier.
	Parent *string `mapstructure:"parent"`

	// Links is a list of task identifiers this task is linked to.
	Links []*string `mapstructure:"links"`

	// ConditionalEffects specifies the effects that can occur based on certain conditions.
	ConditionalEffects []ConditionalEffect `mapstructure:"conditionalEffects"`
}

// To return model.Task
func (t *Task) To() (*model.Task, error) {
	var externalID *domaintask.ExternalID
	var parentID *domaintask.ExternalID
	var links []*domaintask.ExternalID
	var children []model.Task
	var events []domaintask.Event
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
	if t.Parent != nil {
		parentID, err = domaintask.NewExternalID(*t.Parent)
		if err != nil {
			return nil, err
		}
	}
	if t.Links != nil {
		links = make([]*domaintask.ExternalID, len(t.Links))
		for j, link := range t.Links {
			links[j], err = domaintask.NewExternalID(*link)
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
	if t.Events != nil {
		events = make([]domaintask.Event, len(t.Events))
		for i, event := range t.Events {
			d, err := event.To()
			if err != nil {
				return nil, err
			}
			events[i] = *d
		}
	}
	var conditionalDefinitions []*domaintask.ConditionalDefinition
	for _, effect := range t.ConditionalEffects {
		def, err := effect.To()
		if err != nil {
			return nil, err
		}
		conditionalDefinitions = append(conditionalDefinitions, def)
	}

	return &model.Task{
		Name:                  t.Name,
		ExternalID:            externalID,
		Delay:                 *delay,
		Duration:              *duration,
		Kind:                  t.Kind,
		Attributes:            t.Attributes,
		Children:              children,
		ChildOf:               parentID,
		LinkedTo:              links,
		Events:                events,
		ConditionalDefinition: conditionalDefinitions,
	}, nil
}
