package span

import "time"

// Event represents an event in a span
type Event struct {
	name       string
	occurredAt time.Time
	attributes map[string]string
}

func NewEvent(name string, occurredAt time.Time, attributes map[string]string) Event {
	return Event{
		name:       name,
		occurredAt: occurredAt,
		attributes: attributes,
	}
}

// ShiftOccurredAt shifts the occurredAt time by the given offset
func (e *Event) ShiftOccurredAt(offset time.Duration) {
	e.occurredAt = e.occurredAt.Add(offset)
}

// Name returns the name of the event
func (e *Event) Name() string {
	return e.name
}

// OccurredAt returns the time when the event occurred
func (e *Event) OccurredAt() time.Time {
	return e.occurredAt
}

// Attributes returns the attributes of the event
func (e *Event) Attributes() map[string]string {
	return e.attributes
}
