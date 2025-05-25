package task

// Definition represents a task in the trace
type Definition struct {
	name                   string
	isResourceEntryPoint   bool
	resource               *Resource
	attributes             map[string]string
	kind                   Kind
	externalID             *ExternalID
	delay                  Delay                    // Relative time from the start of the parent task
	duration               Duration                 // Relative time from the start of the parent task
	childOf                *ExternalID              // ID of the parent task (if any)
	linkedTo               []*ExternalID            // IDs of linked spans (for producer/consumer relationships)
	events                 []Event                  // Events associated with the task
	conditionalDefinitions []*ConditionalDefinition // Conditional definitions for the task
}

// NewDefinition creates a new task definition
func NewDefinition(name string, isResourceEntryPoint bool, resource *Resource, attributes map[string]string, kind Kind, externalID *ExternalID, delay Delay, duration Duration, childOf *ExternalID, linkedTo []*ExternalID, events []Event, conditionaldefinitions []*ConditionalDefinition) (*Definition, error) {
	return &Definition{
		name:                   name,
		isResourceEntryPoint:   isResourceEntryPoint,
		resource:               resource,
		attributes:             attributes,
		kind:                   kind,
		externalID:             externalID,
		delay:                  delay,
		duration:               duration,
		childOf:                childOf,
		linkedTo:               linkedTo,
		events:                 events,
		conditionalDefinitions: conditionaldefinitions,
	}, nil
}

func (d *Definition) Name() string {
	return d.name
}

func (d *Definition) Resource() *Resource {
	return d.resource
}

func (d *Definition) Attributes() map[string]string {
	return d.attributes
}

func (d *Definition) IsResourceEntryPoint() bool {
	return d.isResourceEntryPoint
}

func (d *Definition) Kind() Kind {
	return d.kind
}

func (d *Definition) ExternalID() *ExternalID {
	return d.externalID
}

func (d *Definition) Delay() Delay {
	return d.delay
}

func (d *Definition) Duration() Duration {
	return d.duration
}

func (d *Definition) ChildOf() *ExternalID {
	return d.childOf
}

func (d *Definition) LinkedTo() []*ExternalID {
	return d.linkedTo
}

func (d *Definition) Events() []Event {
	return d.events
}

func (d *Definition) ConditionalDefinitions() []*ConditionalDefinition {
	return d.conditionalDefinitions
}
