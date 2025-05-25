package span

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

var _ Effect = (*RecordEventEffect)(nil)

// RecordEventEffect is an effect that records an event at a specific time.
type RecordEventEffect struct {
	event task.Event
}

func (r *RecordEventEffect) Apply(node *TreeNode) error {
	duration := node.endTime.Sub(node.startTime)
	if duration < 0 {
		return fmt.Errorf("invalid duration: %v", duration)
	}
	delay, err := r.event.Delay().Resolve(&duration)
	if err != nil {
		return fmt.Errorf("failed to resolve delay: %w", err)
	}
	e := NewEvent(r.event.Name(), node.startTime.Add(*delay), r.event.Attributes())
	node.events = append(node.events, e)
	return nil
}

func FromRecordEventEffect(spec task.RecordEventEffect) (Effect, error) {
	return &RecordEventEffect{event: spec.Event()}, nil
}
