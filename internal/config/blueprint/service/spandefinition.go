package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
	domaintask "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// SpanDefinition represents a span definition in the blueprint.
type SpanDefinition struct {
	// Name is the name of the span.
	Name string `mapstructure:"name"`

	// Ref is an optional identifier for the span.
	Ref *string `mapstructure:"ref"`

	// Delay specifies the delay in duration or relative duration to parent duration before the span starts.
	Delay *Delay `mapstructure:"delay"`

	// Duration specifies the duration of the span.
	Duration *Duration `mapstructure:"duration"`

	// Kind specifies the type or category of the span (e.g., "client", "server").
	Kind string `mapstructure:"kind"`

	// Attributes contains optional attributes for the span.
	Attributes map[string]string `mapstructure:"attributes"`

	// Children is a list of child spans triggered by this span.
	Children []SpanDefinition `mapstructure:"children"`

	// Events is a list of events associated with the span.
	Events []Event `mapstructure:"events"`

	// Parent is an optional parent span ref.
	Parent *string `mapstructure:"parent"`

	// Links is a list of span refs this span is linked to.
	Links []*string `mapstructure:"links"`

	// ConditionalEffects specifies the effects that can occur based on certain conditions.
	ConditionalEffects []ConditionalEffect `mapstructure:"conditional_effects"`
}

// To return model.Task
func (t *SpanDefinition) To() (*model.Task, error) {
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
	if t.Ref != nil {
		externalID, err = domaintask.NewExternalID(*t.Ref)
		if err != nil {
			return nil, fmt.Errorf("invalid external ref: %s", *t.Ref)
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
	var conditionalDefinitions []domaintask.ConditionalDefinition
	for _, effect := range t.ConditionalEffects {
		def, err := effect.To()
		if err != nil {
			return nil, err
		}
		conditionalDefinitions = append(conditionalDefinitions, *def)
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
