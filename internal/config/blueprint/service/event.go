package service

import "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"

type Event struct {
	Name       string            `mapstructure:"name"`
	Delay      Delay             `mapstructure:"delay"`
	Attributes map[string]string `mapstructure:"attributes"`
}

// To converts the event to a domain model.
func (e *Event) To() (*task.Event, error) {
	delay, err := e.Delay.To()
	if err != nil {
		return nil, err
	}

	event := task.NewEvent(e.Name, *delay, e.Attributes)
	return &event, nil
}
