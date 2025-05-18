package service

import (
	"fmt"
	"github.com/k4ji/tracesimulator/pkg/model/task"
)

// Effect represents an effect that can be applied to a span.
type Effect struct {
	// Kind is the kind of effect to be applied.
	Kind string `mapstructure:"kind"`
	// MarkAsFailed is the effect to mark the span as failed.
	MarkAsFailed MarkAsFailed `mapstructure:"markAsFailed"`
	// Annotate is the effect to annotate the span.
	Annotate Annotate `mapstructure:"annotate"`
	// RecordEvent is the effect to record an event.
	RecordEvent RecordEvent `mapstructure:"recordEvent"`
}

// MarkAsFailed represents an effect that marks a span as failed.
type MarkAsFailed struct {
	// Message is the message to be used when marking the span as failed.
	Message string `mapstructure:"message"`
}

// Annotate represents an effect that annotates a span.
type Annotate struct {
	// Attributes are the attributes to be used when annotating the span.
	Attributes map[string]string `mapstructure:"attributes"`
}

// RecordEvent represents an effect that records an event.
type RecordEvent struct {
	Event Event `mapstructure:"event"`
}

// To converts the effect to a domain model.
func (e *Effect) To() (*task.Effect, error) {
	switch e.Kind {
	case "markAsFailed":
		e := task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect(&e.MarkAsFailed.Message))
		return &e, nil
	case "annotate":
		e := task.FromAnnotateEffect(task.NewAnnotateEffect(e.Annotate.Attributes))
		return &e, nil
	case "recordEvent":
		event, err := e.RecordEvent.Event.To()
		if err != nil {
			return nil, fmt.Errorf("failed to convert record event effect: %w", err)
		}
		e := task.FromRecordEventEffect(task.NewRecordEventEffect(*event))
		return &e, nil
	default:
		return nil, fmt.Errorf("unknown effect type: %s", e.Kind)
	}
}
