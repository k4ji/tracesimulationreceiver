package task

// Resource represents an entity that emits spans
type Resource struct {
	name       string            // Name of the resource
	attributes map[string]string // Attributes of the resource
}

// NewResource creates a new Resource with the given name and attributes
func NewResource(name string, attributes map[string]string) *Resource {
	return &Resource{
		name:       name,
		attributes: attributes,
	}
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) Attributes() map[string]string {
	return r.attributes
}
