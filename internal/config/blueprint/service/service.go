package service

import "github.com/k4ji/tracesimulator/pkg/blueprint/service/model"

// Service represents a service in the blueprint.
type Service struct {
	// Name is the name of the service.
	Name string `mapstructure:"name"`

	// Resource contains metadata or attributes associated with the service.
	Resource map[string]string `mapstructure:"resource"`

	// SpanDefinitions is a list of span definitions associated with the service.
	SpanDefinitions []SpanDefinition `mapstructure:"spans"`
}

// To converts the service to a domain model.
func (s *Service) To() (*model.Service, error) {
	tasks := make([]model.Task, 0, len(s.SpanDefinitions))
	for _, sd := range s.SpanDefinitions {
		t, err := sd.To()
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	service := model.Service{
		Name:     s.Name,
		Resource: s.Resource,
		Tasks:    tasks,
	}
	return &service, nil
}
