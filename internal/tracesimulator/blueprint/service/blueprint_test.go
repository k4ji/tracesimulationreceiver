package service

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"github.com/stretchr/testify/assert"
	mathRand "math/rand"
	"testing"
	"time"
)

func TestBlueprint_Interpret(t *testing.T) {
	t.Run("convert each service to trees of tasks", func(t *testing.T) {
		services := []model.Service{
			{
				Name: "service-a",
				Resource: map[string]string{
					"env": "test",
				},
				Tasks: []model.Task{
					{
						Name:       "task-a1",
						ExternalID: func() *task.ExternalID { id, _ := task.NewExternalID("task-a1"); return id }(),
						Delay:      NewAbsoluteDurationDelay(0),
						Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
						Kind:       "server",
						Attributes: map[string]string{
							"key1": "value1",
						},
						Children: []model.Task{
							{
								Name:       "task-a1-child",
								ExternalID: func() *task.ExternalID { id, _ := task.NewExternalID("task-a1-child"); return id }(),
								Delay:      NewAbsoluteDurationDelay(time.Duration(500) * time.Millisecond),
								Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
								Kind:       "producer",
								Attributes: map[string]string{
									"key2": "value2",
								},
								ConditionalDefinition: []task.ConditionalDefinition{
									task.NewConditionalDefinition(
										task.NewProbabilisticCondition(0.1, mathRand.Float64),
										[]task.Effect{
											task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect(ptrString("error"))),
										},
									),
								},
							},
						},
					},
					{
						Name:     "task-a2",
						Kind:     "internal",
						Delay:    NewAbsoluteDurationDelay(0),
						Duration: NewAbsoluteDurationDuration(100 * time.Millisecond),
					},
				},
			},
			{
				Name: "service-b",
				Tasks: []model.Task{
					{Name: "task-b1",
						Kind:     "internal",
						Delay:    NewAbsoluteDurationDelay(0),
						Duration: NewAbsoluteDurationDuration(2000 * time.Millisecond),
					},
				},
			},
		}

		blueprint := NewServiceBlueprint(services)

		rootTaskNodes, err := blueprint.Interpret()

		assert.NoError(t, err)
		assert.Len(t, rootTaskNodes, 3)

		taskA1Resource := rootTaskNodes[0].Definition().Resource()
		assert.Equal(t, "task-a1", rootTaskNodes[0].Definition().Name())
		assert.Equal(t, "service-a", taskA1Resource.Name())
		assert.Equal(t, "test", taskA1Resource.Attributes()["env"])
		assert.Equal(t, task.KindServer, rootTaskNodes[0].Definition().Kind())
		assert.Equal(t, map[string]string{"key1": "value1"}, rootTaskNodes[0].Definition().Attributes())
		assert.Equal(t, NewAbsoluteDurationDelay(0), rootTaskNodes[0].Definition().Delay())
		assert.Equal(t, NewAbsoluteDurationDuration(time.Duration(1000)*time.Millisecond), rootTaskNodes[0].Definition().Duration())
		assert.Len(t, rootTaskNodes[0].Definition().ConditionalDefinitions(), 0)

		taskA1ChildResource := rootTaskNodes[0].Children()[0].Definition().Resource()
		assert.Equal(t, "task-a1-child", rootTaskNodes[0].Children()[0].Definition().Name())
		assert.Equal(t, "service-a", taskA1ChildResource.Name())
		assert.Equal(t, "test", taskA1ChildResource.Attributes()["env"])
		assert.Equal(t, task.KindProducer, rootTaskNodes[0].Children()[0].Definition().Kind())
		assert.Equal(t, map[string]string{"key2": "value2"}, rootTaskNodes[0].Children()[0].Definition().Attributes())
		assert.Equal(t, NewAbsoluteDurationDelay(time.Duration(500)*time.Millisecond), rootTaskNodes[0].Children()[0].Definition().Delay())
		assert.Equal(t, NewAbsoluteDurationDuration(time.Duration(500)*time.Millisecond), rootTaskNodes[0].Children()[0].Definition().Duration())
		assert.Equal(t, 0.1, rootTaskNodes[0].Children()[0].Definition().ConditionalDefinitions()[0].Condition().Probabilistic().Threshold())
		assert.NotNil(t, rootTaskNodes[0].Children()[0].Definition().ConditionalDefinitions()[0].Effects()[0].MarkAsFailedEffect())

		taskA2Resource := rootTaskNodes[1].Definition().Resource()
		assert.Equal(t, "task-a2", rootTaskNodes[1].Definition().Name())
		assert.Equal(t, "service-a", taskA2Resource.Name())
		assert.Equal(t, "test", taskA2Resource.Attributes()["env"])
		assert.Equal(t, task.KindInternal, rootTaskNodes[1].Definition().Kind())
		assert.Equal(t, NewAbsoluteDurationDelay(0), rootTaskNodes[1].Definition().Delay())
		assert.Equal(t, NewAbsoluteDurationDuration(time.Duration(100)*time.Millisecond), rootTaskNodes[1].Definition().Duration())
		assert.Len(t, rootTaskNodes[1].Definition().ConditionalDefinitions(), 0)

		taskB1Resource := rootTaskNodes[2].Definition().Resource()
		assert.Equal(t, "task-b1", rootTaskNodes[2].Definition().Name())
		assert.Equal(t, "service-b", taskB1Resource.Name())
		assert.Equal(t, task.KindInternal, rootTaskNodes[2].Definition().Kind())
		assert.Equal(t, NewAbsoluteDurationDelay(0), rootTaskNodes[2].Definition().Delay())
		assert.Equal(t, NewAbsoluteDurationDuration(time.Duration(2000)*time.Millisecond), rootTaskNodes[2].Definition().Duration())
		assert.Len(t, rootTaskNodes[2].Definition().ConditionalDefinitions(), 0)
	})

	t.Run("connect task nodes across services based on parent-child relationships", func(t *testing.T) {
		parentID, _ := task.NewExternalID("parent-id")
		services := []model.Service{
			{
				Name: "service-a",
				Tasks: []model.Task{
					{
						Name:       "parent-task",
						ExternalID: parentID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
			{
				Name: "service-b",
				Tasks: []model.Task{
					{
						Name:    "child-task",
						Delay:   NewAbsoluteDurationDelay(0),
						ChildOf: parentID,
					},
				},
			},
		}
		blueprint := NewServiceBlueprint(services)

		rootTaskNodes, err := blueprint.Interpret()

		assert.NoError(t, err)
		assert.Len(t, rootTaskNodes, 1)
		assert.Equal(t, "parent-task", rootTaskNodes[0].Definition().Name())
		assert.Len(t, rootTaskNodes[0].Children(), 1)
		assert.Equal(t, "child-task", rootTaskNodes[0].Children()[0].Definition().Name())
	})

	t.Run("connect a child task node and a root task node across services based on parent-child relationships", func(t *testing.T) {
		parentID, _ := task.NewExternalID("parent-id")
		services :=
			[]model.Service{
				{
					Name: "service-a",
					Tasks: []model.Task{
						{
							Name:  "parent-task",
							Delay: NewAbsoluteDurationDelay(0),
							Children: []model.Task{
								{
									Name:       "child-task",
									ExternalID: parentID,
									Delay:      NewAbsoluteDurationDelay(0),
								},
							},
						},
					},
				},
				{
					Name: "service-b",
					Tasks: []model.Task{
						{
							Name:    "child-task",
							Delay:   NewAbsoluteDurationDelay(0),
							ChildOf: parentID,
						},
					},
				},
			}

		blueprint := NewServiceBlueprint(services)

		rootTaskNodes, err := blueprint.Interpret()

		assert.NoError(t, err)
		assert.Len(t, rootTaskNodes, 1)
		assert.Equal(t, "parent-task", rootTaskNodes[0].Definition().Name())
		assert.Len(t, rootTaskNodes[0].Children(), 1)
		assert.Equal(t, "child-task", rootTaskNodes[0].Children()[0].Definition().Name())
	})

	t.Run("link tasks across services based on linkedTo relationships", func(t *testing.T) {
		taskAID, _ := task.NewExternalID("task-a")
		taskBID, _ := task.NewExternalID("task-b")

		services := []model.Service{
			{
				Name: "service-a",
				Tasks: []model.Task{
					{
						Name:       "task-a",
						ExternalID: taskAID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
			{
				Name: "service-b",
				Tasks: []model.Task{
					{
						Name:       "task-b",
						ExternalID: taskBID,
						Delay:      NewAbsoluteDurationDelay(0),
						LinkedTo:   []*task.ExternalID{taskAID},
					},
				},
			},
		}
		blueprint := NewServiceBlueprint(services)

		rootTaskNodes, err := blueprint.Interpret()

		assert.NoError(t, err)
		assert.Len(t, rootTaskNodes, 2)

		taskANode := rootTaskNodes[0]
		taskBNode := rootTaskNodes[1]

		assert.Equal(t, "task-a", taskANode.Definition().Name())
		assert.Equal(t, "task-b", taskBNode.Definition().Name())

		assert.Len(t, taskBNode.Definition().LinkedTo(), 1)
		assert.Equal(t, taskANode.Definition().ExternalID(), taskBNode.Definition().LinkedTo()[0])
	})

	t.Run("return error if the specified parent task is not found", func(t *testing.T) {
		parentID, _ := task.NewExternalID("non-existent-parent")

		services := []model.Service{
			{
				Name: "service-a",
				Tasks: []model.Task{
					{
						Name:    "child-task",
						ChildOf: parentID,
						Delay:   NewAbsoluteDurationDelay(0),
					},
				},
			},
		}
		blueprint := NewServiceBlueprint(services)

		_, err := blueprint.Interpret()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent task not found")
	})

	t.Run("return error if duplicate ExternalID is detected", func(t *testing.T) {
		externalID, _ := task.NewExternalID("duplicate-id")

		services := []model.Service{
			{
				Name: "service-a",
				Tasks: []model.Task{
					{
						Name:       "task-a1",
						ExternalID: externalID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
			{
				Name: "service-b",
				Tasks: []model.Task{
					{
						Name:       "task-b1",
						ExternalID: externalID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
		}
		blueprint := NewServiceBlueprint(services)

		_, err := blueprint.Interpret()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate ExternalID detected")
	})

	t.Run("return error if cyclic dependencies are detected", func(t *testing.T) {
		taskAID, _ := task.NewExternalID("task-a")
		taskBID, _ := task.NewExternalID("task-b")

		services := []model.Service{
			{
				Name: "service-a",
				Tasks: []model.Task{
					{
						Name:       "task-a",
						ExternalID: taskAID,
						ChildOf:    taskBID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
			{
				Name: "service-b",
				Tasks: []model.Task{
					{
						Name:       "task-b",
						ExternalID: taskBID,
						ChildOf:    taskAID,
						Delay:      NewAbsoluteDurationDelay(0),
					},
				},
			},
		}
		blueprint := NewServiceBlueprint(services)

		_, err := blueprint.Interpret()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to add child task-b to parent task-a: adding task-b as child of task-a would create a cycle")
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
