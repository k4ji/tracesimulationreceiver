package task

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"testing"
	"time"
)

func makeDefaultDefinition(name string) *Definition {
	def, _ := NewDefinition(
		name,
		false,
		NewResource("test_service", make(map[string]string)),
		make(map[string]string),
		KindInternal,
		nil,
		NewAbsoluteDurationDelay(0),
		NewAbsoluteDurationDuration(0),
		nil,
		make([]*ExternalID, 0),
		make([]Event, 0),
		[]*ConditionalDefinition{
			NewConditionalDefinition(
				Condition{
					kind:          ConditionKindProbabilistic,
					probabilistic: &ProbabilisticCondition{threshold: 0.0},
				},
				make([]Effect, 0),
			),
		},
	)
	return def
}

func TestTreeNode_AddChild(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (*TreeNode, *TreeNode, *TreeNode)
		action    func(a, b, c *TreeNode) error
		expectErr bool
	}{
		{
			name: "success",
			setup: func() (*TreeNode, *TreeNode, *TreeNode) {
				a := NewTreeNode(makeDefaultDefinition("A"))
				b := NewTreeNode(makeDefaultDefinition("B"))
				c := NewTreeNode(makeDefaultDefinition("C"))
				return a, b, c
			},
			action: func(a, b, c *TreeNode) error {
				if err := a.AddChild(b); err != nil {
					return err
				}
				return b.AddChild(c)
			},
			expectErr: false,
		},
		{
			name: "return error when cyclic dependency is detected",
			setup: func() (*TreeNode, *TreeNode, *TreeNode) {
				a := NewTreeNode(makeDefaultDefinition("A"))
				b := NewTreeNode(makeDefaultDefinition("B"))
				c := NewTreeNode(makeDefaultDefinition("C"))
				_ = a.AddChild(b)
				_ = b.AddChild(c)
				return a, b, c
			},
			action: func(a, b, c *TreeNode) error {
				return c.AddChild(a)
			},
			expectErr: true,
		},
		{
			name: "return error when duplicate parent is detected",
			setup: func() (*TreeNode, *TreeNode, *TreeNode) {
				a := NewTreeNode(makeDefaultDefinition("A"))
				b := NewTreeNode(makeDefaultDefinition("B"))
				c := NewTreeNode(makeDefaultDefinition("C"))
				_ = a.AddChild(c)
				return a, b, c
			},
			action: func(a, b, c *TreeNode) error {
				return b.AddChild(c)
			},
			expectErr: true,
		},
		{
			name: "return error when self parent is detected",
			setup: func() (*TreeNode, *TreeNode, *TreeNode) {
				a := NewTreeNode(makeDefaultDefinition("A"))
				return a, nil, nil
			},
			action: func(a, b, c *TreeNode) error {
				return a.AddChild(a)
			},
			expectErr: true,
		},
		{
			name: "return error when nil child is detected",
			setup: func() (*TreeNode, *TreeNode, *TreeNode) {
				a := NewTreeNode(makeDefaultDefinition("A"))
				return a, nil, nil
			},
			action: func(a, b, c *TreeNode) error {
				return a.AddChild(nil)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b, c := tt.setup()
			err := tt.action(a, b, c)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func NewAbsoluteDurationDelay(duration time.Duration) Delay {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := NewDelay(e)
	return *d
}

func NewAbsoluteDurationDuration(duration time.Duration) Duration {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := NewDuration(e)
	return *d
}
