package simulator

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/adapter"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/span"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"github.com/stretchr/testify/assert"
	mathRand "math/rand"
	"testing"
	"time"
)

type MockAdapter struct{}

var _ simulator.Adapter[[]string] = &MockAdapter{}

func (m *MockAdapter) Transform(spans []*span.TreeNode) ([]string, error) {
	names := make([]string, len(spans))
	for i, s := range spans {
		names[i] = s.Name()
	}
	return names, nil
}

func TestSimulator_Run(t *testing.T) {
	rootTaskAExternalID, _ := task.NewExternalID("root-a")
	childTaskA1ExternalID, _ := task.NewExternalID("child-a1")
	childTaskA2ExternalID, _ := task.NewExternalID("child-a2")
	rootTaskBExternalID, _ := task.NewExternalID("root-b")
	rootTaskCExternalID, _ := task.NewExternalID("root-c")
	now := time.Now()

	blueprint := service.NewServiceBlueprint([]model.Service{
		{
			Name: "service-a",
			Resource: map[string]string{
				"env": "test",
			},
			Tasks: []model.Task{
				{
					Name:       "root-task-a",
					ExternalID: rootTaskAExternalID,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
					Kind:       "server",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-a-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-a-1": "attribute-value-event-root-task-a-1",
							},
						),
						task.NewEvent(
							"event-root-task-a-2",
							NewAbsoluteDurationDelay(100*time.Millisecond),
							map[string]string{
								"attribute-key-event-root-task-a-2": "attribute-value-event-root-task-a-2",
							},
						),
					},
					Attributes: map[string]string{
						"key1": "value1",
					},
					Children: []model.Task{
						{
							Name:       "child-task-a1",
							ExternalID: childTaskA1ExternalID,
							Delay:      NewAbsoluteDurationDelay(time.Duration(500) * time.Millisecond),
							Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
							Kind:       "producer",
							Events: []task.Event{
								task.NewEvent(
									"event-child-task-a1-1",
									NewAbsoluteDurationDelay(0),
									map[string]string{
										"attribute-key-event-child-task-a1-1": "attribute-value-event-child-task-a1-1",
									},
								),
							},
						},
						{
							Name:       "child-task-a2",
							ExternalID: childTaskA2ExternalID,
							Delay:      NewAbsoluteDurationDelay(time.Duration(1000) * time.Millisecond),
							Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
							Kind:       "internal",
							Events: []task.Event{
								task.NewEvent(
									"event-child-task-a2-1",
									NewAbsoluteDurationDelay(0),
									map[string]string{
										"attribute-key-event-child-task-a2-1": "attribute-value-event-child-task-a2-1",
									},
								),
							},
						},
					},
				},
			},
		},
		{
			Name: "service-b",
			Tasks: []model.Task{
				{
					Name:       "root-task-b",
					ExternalID: rootTaskBExternalID,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
					Kind:       "consumer",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-b-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-b-1": "attribute-value-event-root-task-b-1",
							},
						),
					},
					Children: []model.Task{
						{
							Name:       "child-task-b1",
							ExternalID: nil,
							Delay:      NewAbsoluteDurationDelay(time.Duration(1000) * time.Millisecond),
							Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
							Kind:       "client",
							Events: []task.Event{
								task.NewEvent(
									"event-child-task-b1-1",
									NewAbsoluteDurationDelay(0),
									map[string]string{
										"attribute-key-event-child-task-b1-1": "attribute-value-event-child-task-b1-1",
									},
								),
							},
						},
					},
					LinkedTo: []*task.ExternalID{
						childTaskA1ExternalID,
					},
				},
			},
		},
		// spans with childOf
		{
			Name: "service-c",
			Tasks: []model.Task{
				{
					Name:       "root-task-c",
					ExternalID: rootTaskCExternalID,
					Delay:      NewAbsoluteDurationDelay(time.Duration(1000) * time.Millisecond),
					Duration:   NewAbsoluteDurationDuration(3000 * time.Millisecond),
					Kind:       "internal",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-c-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-c-1": "attribute-value-event-root-task-c-1",
							},
						),
					},
					ChildOf: childTaskA2ExternalID,
				},
			},
		},
		// spans with error probability
		{
			Name: "service-d",
			Tasks: []model.Task{
				{
					Name:       "root-task-d",
					ExternalID: nil,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(1500 * time.Millisecond),
					Kind:       "internal",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-d-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-d-1": "attribute-value-event-root-task-d-1",
							},
						),
					},
					ChildOf: nil,
					ConditionalDefinition: []task.ConditionalDefinition{
						task.NewConditionalDefinition(
							task.NewProbabilisticCondition(1.0, mathRand.Float64),
							[]task.Effect{
								task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
							},
						),
					},
				},
			},
		},
	})

	sim := New[[]*span.TreeNode](&simulator.NoOpAdapter{})
	traces, err := sim.Run(&blueprint, now)
	assert.NoError(t, err)

	t.Run("create spans by populating the blueprint", func(t *testing.T) {
		assert.Len(t, traces, 3)

		// Validate the first trace (service-a)
		rootA := traces[0]
		traceID := rootA.TraceID()
		rootAResource := rootA.Resource()
		assert.Equal(t, "root-task-a", rootA.Name())
		assert.Equal(t, "service-a", rootAResource.Name())
		assert.Equal(t, true, rootA.IsResourceEntryPoint())
		assert.Equal(t, map[string]string{"env": "test"}, rootAResource.Attributes())
		assert.Equal(t, map[string]string{"key1": "value1"}, rootA.Attributes())
		assert.Equal(t, span.KindServer, rootA.Kind())
		assert.Equal(t, rootTaskAExternalID, rootA.ExternalID())

		assert.Len(t, rootA.Children(), 2)

		childA1 := rootA.Children()[0]
		childA1Resource := childA1.Resource()
		assert.Equal(t, traceID, childA1.TraceID())
		assert.Equal(t, "child-task-a1", childA1.Name())
		assert.Equal(t, "service-a", childA1Resource.Name())
		assert.Equal(t, false, childA1.IsResourceEntryPoint())
		assert.Equal(t, span.KindProducer, childA1.Kind())
		assert.Equal(t, childTaskA1ExternalID, childA1.ExternalID())

		childA2 := rootA.Children()[1]
		childA2Resource := childA2.Resource()
		assert.Equal(t, traceID, childA2.TraceID())
		assert.Equal(t, "child-task-a2", childA2.Name())
		assert.Equal(t, "service-a", childA2Resource.Name())
		assert.Equal(t, false, childA2.IsResourceEntryPoint())
		assert.Equal(t, span.KindInternal, childA2.Kind())
		assert.Len(t, childA2.Children(), 1)

		// Validate the child task, which is a service entry point
		rootC := childA2.Children()[0]
		rootCResource := rootC.Resource()
		assert.Equal(t, traceID, rootC.TraceID())
		assert.Equal(t, "root-task-c", rootC.Name())
		assert.Equal(t, "service-c", rootCResource.Name())
		assert.Equal(t, true, rootC.IsResourceEntryPoint())
		assert.Equal(t, span.KindInternal, rootC.Kind())
		assert.Equal(t, rootTaskCExternalID, rootC.ExternalID())

		// Validate the second trace (service-b)
		rootB := traces[1]
		rootBResource := rootB.Resource()
		assert.Equal(t, "root-task-b", rootB.Name())
		assert.Equal(t, "service-b", rootBResource.Name())
		assert.Equal(t, true, rootB.IsResourceEntryPoint())
		assert.Equal(t, span.KindConsumer, rootB.Kind())
		assert.Equal(t, rootTaskBExternalID, rootB.ExternalID())
		assert.Len(t, rootB.Children(), 1)

		childB1 := rootB.Children()[0]
		childB1Resource := childB1.Resource()
		assert.Equal(t, "child-task-b1", childB1.Name())
		assert.Equal(t, "service-b", childB1Resource.Name())
		assert.Equal(t, false, childB1.IsResourceEntryPoint())
		assert.Equal(t, span.KindClient, childB1.Kind())
		assert.Nil(t, childB1.ExternalID())

		// Validate the third trace (service-d)
		rootD := traces[2]
		rootDResource := rootD.Resource()
		assert.Equal(t, "root-task-d", rootD.Name())
		assert.Equal(t, "service-d", rootDResource.Name())
		assert.Equal(t, true, rootD.IsResourceEntryPoint())
		assert.Equal(t, span.KindInternal, rootD.Kind())
		assert.Nil(t, rootD.ExternalID())
	})

	t.Run("link spans via external IDs", func(t *testing.T) {
		rootB := traces[1]
		assert.Len(t, rootB.LinkedTo(), 1)
		linkedSpan := rootB.LinkedTo()[0]
		assert.Equal(t, "child-task-a1", linkedSpan.Name())
		linkedSpanResource := linkedSpan.Resource()
		assert.Equal(t, "service-a", linkedSpanResource.Name())
	})

	t.Run("adjust span timestamps to ensure all end before base end time", func(t *testing.T) {
		// The longest trace takes 5000ms where childA1, which has startAfter of 1000ms, plus rootC, which has startAfter of 1000ms and duration of 3000ms.
		expectedStartTime := now.Add(-5000 * time.Millisecond)
		rootA := traces[0]
		assert.GreaterOrEqual(t, rootA.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, rootA.EndTime(), now)

		childA1 := rootA.Children()[0]
		assert.GreaterOrEqual(t, childA1.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, childA1.EndTime(), now)

		rootC := rootA.Children()[1]
		assert.GreaterOrEqual(t, rootC.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, rootC.EndTime(), now)

		rootB := traces[1]
		assert.GreaterOrEqual(t, rootB.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, rootB.EndTime(), now)

		childB1 := rootB.Children()[0]
		assert.GreaterOrEqual(t, childB1.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, childB1.EndTime(), now)

		rootD := traces[2]
		assert.GreaterOrEqual(t, rootD.StartTime(), expectedStartTime)
		assert.LessOrEqual(t, rootD.EndTime(), now)
	})

	t.Run("adjust span event timestamps to ensure all events are within the span's time range", func(t *testing.T) {
		rootA := traces[0]
		assert.Len(t, rootA.Events(), 2)
		assert.GreaterOrEqual(t, rootA.Events()[0].OccurredAt(), rootA.StartTime())
		assert.LessOrEqual(t, rootA.Events()[0].OccurredAt(), now)
		assert.GreaterOrEqual(t, rootA.Events()[1].OccurredAt(), rootA.StartTime())
		assert.LessOrEqual(t, rootA.Events()[1].OccurredAt(), now)

		childA1 := rootA.Children()[0]
		assert.Len(t, childA1.Events(), 1)
		assert.GreaterOrEqual(t, childA1.Events()[0].OccurredAt(), childA1.StartTime())
		assert.LessOrEqual(t, childA1.Events()[0].OccurredAt(), now)

		childA2 := rootA.Children()[1]
		assert.Len(t, childA2.Events(), 1)
		assert.GreaterOrEqual(t, childA2.Events()[0].OccurredAt(), childA2.StartTime())
		assert.LessOrEqual(t, childA2.Events()[0].OccurredAt(), now)

		rootC := childA2.Children()[0]
		assert.Len(t, rootC.Events(), 1)
		assert.GreaterOrEqual(t, rootC.Events()[0].OccurredAt(), rootC.StartTime())
		assert.LessOrEqual(t, rootC.Events()[0].OccurredAt(), now)

		rootB := traces[1]
		assert.Len(t, rootB.Events(), 1)
		assert.GreaterOrEqual(t, rootB.Events()[0].OccurredAt(), rootB.StartTime())
		assert.LessOrEqual(t, rootB.Events()[0].OccurredAt(), now)

		childB1 := rootB.Children()[0]
		assert.Len(t, childB1.Events(), 1)
		assert.GreaterOrEqual(t, childB1.Events()[0].OccurredAt(), childB1.StartTime())
		assert.LessOrEqual(t, childB1.Events()[0].OccurredAt(), now)

		rootD := traces[2]
		assert.Len(t, rootD.Events(), 1)
		assert.GreaterOrEqual(t, rootD.Events()[0].OccurredAt(), rootD.StartTime())
		assert.LessOrEqual(t, rootD.Events()[0].OccurredAt(), now)
	})

	t.Run("create error spans based on the probability", func(t *testing.T) {
		assert.Len(t, traces, 3)

		rootA := traces[0]
		assert.Equal(t, span.StatusOK, rootA.Status())

		childA1 := rootA.Children()[0]
		assert.Equal(t, span.StatusOK, childA1.Status())

		rootC := rootA.Children()[1]
		assert.Equal(t, span.StatusOK, rootC.Status())

		rootB := traces[1]
		assert.Equal(t, span.StatusOK, rootB.Status())

		childB1 := rootB.Children()[0]
		assert.Equal(t, span.StatusOK, childB1.Status())

		rootD := traces[2]
		assert.Equal(t, span.StatusError("error"), rootD.Status())
	})

	t.Run("fail to link spans with missing external IDs", func(t *testing.T) {
		missingExternalID, _ := task.NewExternalID("missing-external-id")

		// Create a mock blueprint with linked span without specified external ID
		missingExternalIDBlueprint := service.NewServiceBlueprint([]model.Service{
			{
				Name: "service-x",
				Tasks: []model.Task{
					{
						Name:       "task-without-external-id",
						ExternalID: nil,
						Delay:      NewAbsoluteDurationDelay(0),
						Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
						Kind:       "internal",
					},
					{
						Name:       "task-with-link",
						ExternalID: nil,
						Delay:      NewAbsoluteDurationDelay(0),
						Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
						Kind:       "internal",
						LinkedTo: []*task.ExternalID{
							missingExternalID,
						},
					},
				},
			},
		})

		sim := New[[]*span.TreeNode](&simulator.NoOpAdapter{})
		_, err := sim.Run(&missingExternalIDBlueprint, time.Now())
		assert.Errorf(t, err, "failed to link spans: linked span with external ID {%s} not found", missingExternalID)
	})

	t.Run("fail to construct task tree with duplicate external IDs within the same trace", func(t *testing.T) {
		duplicateExternalID, _ := task.NewExternalID("external-id")

		duplicateExternalIDBlueprint :=
			service.NewServiceBlueprint([]model.Service{
				{
					Name: "service-a",
					Tasks: []model.Task{
						{
							Name:       "root-task-a",
							ExternalID: duplicateExternalID,
							Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
							Kind:       "internal",
							Children: []model.Task{
								{
									Name:       "child-task-a1",
									ExternalID: duplicateExternalID,
									Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
									Kind:       "internal",
								},
							},
						},
					},
				},
			})
		sim := New[[]*span.TreeNode](&simulator.NoOpAdapter{})
		_, err := sim.Run(&duplicateExternalIDBlueprint, time.Now())
		assert.Errorf(t, err, "failed to convert task tree to span:, duplicate external ID {%s}", duplicateExternalID)
	})

	t.Run("fail to construct task tree with duplicate external IDs across multiple traces", func(t *testing.T) {
		duplicateExternalID, _ := task.NewExternalID("external-id")

		duplicateExternalIDBlueprint :=
			service.NewServiceBlueprint([]model.Service{
				{
					Name: "service-a",
					Tasks: []model.Task{
						{
							Name:       "root-task-a",
							ExternalID: duplicateExternalID,
							Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
							Kind:       "internal",
						},
					},
				},
				{
					Name: "service-b",
					Tasks: []model.Task{
						{
							Name:       "root-task-b",
							ExternalID: duplicateExternalID,
							Duration:   NewAbsoluteDurationDuration(2000 * time.Millisecond),
							Kind:       "internal",
						},
					},
				},
			})

		sim := New[[]*span.TreeNode](&simulator.NoOpAdapter{})
		_, err := sim.Run(&duplicateExternalIDBlueprint, time.Now())
		assert.Errorf(t, err, "failed to interpret blueprint: duplicate ExternalID detected: {%s}", duplicateExternalID)
	})

	t.Run("transform span trees to a different format using the adapter", func(t *testing.T) {
		sim := New[[]string](&MockAdapter{})
		transformed, err := sim.Run(&blueprint, time.Now())
		assert.NoError(t, err)
		assert.Len(t, transformed, 3)
		assert.Equal(t, "root-task-a", transformed[0])
		assert.Equal(t, "root-task-b", transformed[1])
		assert.Equal(t, "root-task-d", transformed[2])
	})
}

func NewAbsoluteDurationDelay(duration time.Duration) task.Delay {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := task.NewDelay(e)
	return *d
}

func NewAbsoluteDurationDuration(duration time.Duration) task.Duration {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := task.NewDuration(e)
	return *d
}

func ptrString(s string) *string {
	return &s
}
