package task

// Event represents an event associated with a task
type Event struct {
	name       string
	delay      Delay
	attributes map[string]string
}

func NewEvent(name string, delay Delay, attributes map[string]string) Event {
	return Event{
		name:       name,
		delay:      delay,
		attributes: attributes,
	}
}

// Name returns the name of the event
func (e *Event) Name() string {
	return e.name
}

// Delay returns the delay of the event relative to the task
func (e *Event) Delay() Delay {
	return e.delay
}

// Attributes returns the attributes of the event
func (e *Event) Attributes() map[string]string {
	return e.attributes
}
