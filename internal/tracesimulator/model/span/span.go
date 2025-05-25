package span

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"time"
)

// TreeNode represents a span in the span tree
type TreeNode struct {
	id                   ID
	traceID              TraceID
	name                 string
	isResourceEntryPoint bool
	resource             *task.Resource
	attributes           map[string]string
	kind                 Kind
	startTime            time.Time
	endTime              time.Time
	parentID             *ID
	externalID           *task.ExternalID
	children             []*TreeNode
	linkedTo             []*TreeNode
	events               []Event
	linkedToExternalID   []*task.ExternalID
	status               Status
}

// FromTaskTree converts a task tree to a span tree
func FromTaskTree(
	taskTree *task.TreeNode,
	traceID TraceID,
	baseStartTime time.Time,
	idGen func() ID,
) (*TreeNode, error) {
	rootSpan, err := fromTaskNode(taskTree, traceID, nil, nil, baseStartTime, idGen)
	if err != nil {
		return nil, fmt.Errorf("failed to convert task tree to span tree: %w", err)
	}
	if err := rootSpan.validate(); err != nil {
		return nil, err
	}
	return rootSpan, nil
}

func fromTaskNode(
	taskNode *task.TreeNode,
	traceID TraceID,
	parentID *ID,
	parentDuration *time.Duration,
	baseStartTime time.Time,
	idGen func() ID,
) (*TreeNode, error) {
	spanID := idGen()
	delay, err := taskNode.Definition().Delay().Resolve(parentDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve delay: %w", err)
	}
	duration, err := taskNode.Definition().Duration().Resolve(parentDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve duration: %w", err)
	}

	startTime := baseStartTime.Add(*delay)
	endTime := startTime.Add(*duration)

	events := make([]Event, len(taskNode.Definition().Events()))
	for i, event := range taskNode.Definition().Events() {
		d, err := event.Delay().Resolve(duration)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve event delay: %w", err)
		}
		if *d > *duration {
			return nil, fmt.Errorf("event delay cannot be greater than task duration")
		}
		events[i] = NewEvent(
			event.Name(),
			startTime.Add(*d),
			event.Attributes(),
		)
	}

	node := TreeNode{
		id:                   spanID,
		traceID:              traceID,
		name:                 taskNode.Definition().Name(),
		isResourceEntryPoint: taskNode.Definition().IsResourceEntryPoint(),
		resource:             taskNode.Definition().Resource(),
		attributes:           taskNode.Definition().Attributes(),
		kind:                 FromTaskKind(taskNode.Definition().Kind()),
		startTime:            startTime,
		endTime:              endTime,
		parentID:             parentID,
		externalID:           taskNode.Definition().ExternalID(),
		children:             []*TreeNode{},
		linkedTo:             []*TreeNode{},
		events:               events,
		linkedToExternalID:   taskNode.Definition().LinkedTo(),
		status:               StatusOK,
	}

	for _, childTask := range taskNode.Children() {
		childSpan, err := fromTaskNode(childTask, traceID, &spanID, duration, startTime, idGen)
		if err != nil {
			return nil, fmt.Errorf("failed to convert child task to span: %w", err)
		}
		node.children = append(node.children, childSpan)
	}

	for _, spec := range taskNode.Definition().ConditionalDefinitions() {
		condition, err := FromConditionSpec(spec.Condition())
		if err != nil {
			return nil, fmt.Errorf("failed to convert condition spec to condition: %w", err)
		}
		cr, err := condition.Evaluate(&node)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate condition: %w", err)
		}
		is, err := cr.IsSatisfied()
		if err != nil {
			return nil, fmt.Errorf("failed to check condition satisfaction: %w", err)
		}
		if is {
			for _, effectSpec := range spec.Effects() {
				effect, err := FromEffectSpec(effectSpec)
				if err != nil {
					return nil, fmt.Errorf("failed to convert effect spec to effect: %w", err)
				}
				if err := effect.Apply(&node); err != nil {
					return nil, fmt.Errorf("failed to apply effect: %w", err)
				}
			}
		}
	}

	return &node, nil
}

func (n *TreeNode) validate() error {
	// it returns an error if the externalID is not unique
	externalIDToSpan := make(map[task.ExternalID]*TreeNode)
	var checkDuplicateExternalID func(n *TreeNode) error
	checkDuplicateExternalID = func(n *TreeNode) error {
		for _, child := range n.children {
			err := checkDuplicateExternalID(child)
			if err != nil {
				return err
			}
		}
		if n.externalID != nil {
			if _, exists := externalIDToSpan[*n.externalID]; exists {
				err := fmt.Errorf("duplicate external ID %s", *n.externalID)
				return err
			}
			externalIDToSpan[*n.externalID] = n
		}
		return nil
	}
	return checkDuplicateExternalID(n)
}

// ShiftTimestamps shifts the start and end timestamps of the span and its children by a given duration
func (n *TreeNode) ShiftTimestamps(delta time.Duration) {
	n.startTime = n.startTime.Add(delta)
	n.endTime = n.endTime.Add(delta)
	for i := range n.events {
		n.events[i].ShiftOccurredAt(delta)
	}
	for _, child := range n.children {
		child.ShiftTimestamps(delta)
	}
}

// ExternalIDToSpan returns a map of external IDs to the span and its children
func (n *TreeNode) ExternalIDToSpan() map[task.ExternalID]*TreeNode {
	// it returns an error if the externalID is not unique
	externalIDToSpan := make(map[task.ExternalID]*TreeNode)
	if n.externalID != nil {
		externalIDToSpan[*n.externalID] = n
	}
	for _, child := range n.children {
		childExternalIDToSpan := child.ExternalIDToSpan()
		for id, span := range childExternalIDToSpan {
			externalIDToSpan[id] = span
		}
	}
	return externalIDToSpan
}

// LinkSpan links the spans based on their external IDs and map of external IDs to spans
func (n *TreeNode) LinkSpan(externalIDToSpan map[task.ExternalID]*TreeNode) error {
	for _, linkedToExternalID := range n.linkedToExternalID {
		if linkedSpan, exists := externalIDToSpan[*linkedToExternalID]; exists {
			n.linkedTo = append(n.linkedTo, linkedSpan)
		} else {
			return fmt.Errorf("linked span with external ID %s not found", *linkedToExternalID)
		}
	}
	for _, child := range n.children {
		err := child.LinkSpan(externalIDToSpan)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *TreeNode) ID() ID {
	return n.id
}

func (n *TreeNode) TraceID() TraceID {
	return n.traceID
}

func (n *TreeNode) Name() string {
	return n.name
}

func (n *TreeNode) IsResourceEntryPoint() bool {
	return n.isResourceEntryPoint
}

func (n *TreeNode) Resource() *task.Resource {
	return n.resource
}

func (n *TreeNode) Attributes() map[string]string {
	return n.attributes
}

func (n *TreeNode) Kind() Kind {
	return n.kind
}

func (n *TreeNode) StartTime() time.Time {
	return n.startTime
}

func (n *TreeNode) EndTime() time.Time {
	return n.endTime
}

func (n *TreeNode) ParentID() *ID {
	return n.parentID
}

func (n *TreeNode) ExternalID() *task.ExternalID {
	return n.externalID
}

func (n *TreeNode) Children() []*TreeNode {
	cp := make([]*TreeNode, len(n.children))
	copy(cp, n.children)
	return cp
}

func (n *TreeNode) LinkedTo() []*TreeNode {
	cp := make([]*TreeNode, len(n.linkedTo))
	copy(cp, n.linkedTo)
	return cp
}

func (n *TreeNode) Events() []Event {
	cp := make([]Event, len(n.events))
	copy(cp, n.events)
	return cp
}

func (n *TreeNode) LinkedToExternalID() []*task.ExternalID {
	cp := make([]*task.ExternalID, len(n.linkedToExternalID))
	copy(cp, n.linkedToExternalID)
	return cp
}

func (n *TreeNode) Status() Status {
	return n.status
}
