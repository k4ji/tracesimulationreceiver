package task

type AnnotateEffect struct {
	// attributes is a map of attributes to be added to the task.
	attributes map[string]string
}

// NewAnnotateEffect creates a new AnnotateEffect with the given attributes.
func NewAnnotateEffect(attributes map[string]string) AnnotateEffect {
	return AnnotateEffect{
		attributes: attributes,
	}
}

// Attributes returns the attributes to be added to the task.
func (a AnnotateEffect) Attributes() map[string]string {
	return a.attributes
}
