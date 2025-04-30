package service

import (
	"github.com/k4ji/tracesimulator/pkg/model/task"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	t.Run("valid blueprint", func(t *testing.T) {
		duration := time.Duration(1000)
		bp := &Blueprint{
			Default: DefaultValues{
				StartAfter: &duration,
				Duration:   &duration,
				FailWith: FailureCondition{
					Probability: floatPtr(0.5),
				},
			},
			Services: []Service{},
		}
		err := bp.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing startAfter with default", func(t *testing.T) {
		duration := time.Duration(1000)
		bp := &Blueprint{
			Default: DefaultValues{
				StartAfter: &duration,
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:     "task1",
							Kind:     "client",
							Duration: ptrDuration(1000),
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing startAfter without default", func(t *testing.T) {
		bp := &Blueprint{
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:     "task1",
							Kind:     "client",
							Duration: ptrDuration(1000),
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a StartAfter value or global default")
	})

	t.Run("missing duration with default", func(t *testing.T) {
		duration := time.Duration(1000)
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &duration,
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							StartAfter: ptrDuration(1000),
							Kind:       "client",
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing duration without default", func(t *testing.T) {
		bp := &Blueprint{
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							StartAfter: ptrDuration(1000),
							Kind:       "client",
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a Duration value or global default")
	})

	t.Run("missing failWith with default", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				FailWith: FailureCondition{
					Probability: floatPtr(0.5),
				},
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							StartAfter: ptrDuration(1000),
							Duration:   ptrDuration(1000),
							Kind:       "client",
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing failWith without default", func(t *testing.T) {
		bp := &Blueprint{
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							StartAfter: ptrDuration(1000),
							Duration:   ptrDuration(1000),
							Kind:       "client",
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a FailWithProbability value or global default")
	})

	t.Run("invalid default startAfter", func(t *testing.T) {
		duration := time.Duration(-1)
		bp := &Blueprint{
			Default: DefaultValues{
				StartAfter: &duration,
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "global default wait must be greater than or equal to 0")
	})

	t.Run("invalid default duration", func(t *testing.T) {
		duration := time.Duration(-1)
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &duration,
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "global default duration must be greater than 0")
	})

	t.Run("invalid default failWith probability", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				FailWith: FailureCondition{
					Probability: floatPtr(1.5),
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "global default failWith.probability value between 0.0 and 1.0")
	})
}

func TestTo(t *testing.T) {
	t.Run("convert config to blueprint", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				StartAfter: ptrDuration(100),
				Duration:   ptrDuration(200),
				FailWith: FailureCondition{
					Probability: floatPtr(0.5),
				},
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							ExternalID: ptrString("task1"),
							StartAfter: ptrDuration(1 * time.Millisecond),
							Duration:   ptrDuration(2 * time.Millisecond),
							Kind:       "client",
							Attributes: map[string]string{
								"key1": "value1",
							},
							Children: []Task{
								{
									Name:       "task1-child",
									StartAfter: ptrDuration(3 * time.Millisecond),
									Duration:   ptrDuration(4 * time.Millisecond),
									Kind:       "internal",
									Attributes: map[string]string{
										"key2": "value2",
									},
									FailWith: FailureCondition{
										Probability: floatPtr(0.2),
									},
								},
							},
							FailWith: FailureCondition{
								Probability: floatPtr(0.1),
							},
						},
					},
				},
				{
					Name: "service2",
					Tasks: []Task{
						{
							// childOf plus missing startAfter
							Name:     "task2",
							Duration: ptrDuration(5),
							Kind:     "server",
							Children: []Task{
								// missing Duration
								{
									Name:       "task2-child1",
									ExternalID: ptrString("task2-child1"),
									StartAfter: ptrDuration(6 * time.Millisecond),
									Kind:       "producer",
									FailWith: FailureCondition{
										Probability: floatPtr(0.4),
									},
								},
								// missing FailWith
								{
									Name:       "task2-child2",
									ExternalID: ptrString("task2-child2"),
									StartAfter: ptrDuration(7 * time.Millisecond),
									Duration:   ptrDuration(8 * time.Millisecond),
									Kind:       "producer",
								},
							},
							ChildOf: ptrString("task1"),
							FailWith: FailureCondition{
								Probability: floatPtr(0.3),
							},
						},
					},
				},
				{
					Name: "service3",
					Tasks: []Task{
						// LinkedTo plus missing StartAfter, Duration and FailWith
						{
							Name: "task3",
							Kind: "consumer",
							LinkedTo: []*string{
								ptrString("task2-child1"),
								ptrString("task2-child2"),
							},
						},
					},
				},
			}}
		sbp, err := bp.To()
		assert.NoError(t, err)
		result, err := sbp.Interpret()
		assert.NoError(t, err)
		// compare all fields of the blueprint
		assert.Len(t, result, 2)
		// check the first task tree
		task1 := result[0]
		assert.Equal(t, "task1", task1.Definition().Name())
		assert.Equal(t, "value1", task1.Definition().Attributes()["key1"])
		assert.Equal(t, task.KindClient, task1.Definition().Kind())
		assert.Equal(t, time.Duration(1)*time.Millisecond, task1.Definition().StartAfter())
		assert.Equal(t, time.Duration(2)*time.Millisecond, task1.Definition().Duration())
		assert.Equal(t, 0.1, task1.Definition().FailWithProbability())

		assert.Equal(t, "task1-child", task1.Children()[0].Definition().Name())
		assert.Equal(t, "value2", task1.Children()[0].Definition().Attributes()["key2"])
		assert.Equal(t, task.KindInternal, task1.Children()[0].Definition().Kind())
		assert.Equal(t, time.Duration(3)*time.Millisecond, task1.Children()[0].Definition().StartAfter())
		assert.Equal(t, time.Duration(4)*time.Millisecond, task1.Children()[0].Definition().Duration())
		assert.Equal(t, 0.2, task1.Children()[0].Definition().FailWithProbability())

		task2 := result[0].Children()[1]
		assert.Equal(t, "task2", task2.Definition().Name())
		assert.Equal(t, task.KindServer, task2.Definition().Kind())
		assert.Equal(t, time.Duration(100), task2.Definition().StartAfter())
		assert.Equal(t, time.Duration(5), task2.Definition().Duration())
		assert.Equal(t, "task1", task2.Definition().ChildOf().Value())
		assert.Equal(t, 0.3, task2.Definition().FailWithProbability())

		assert.Equal(t, "task2-child1", task2.Children()[0].Definition().Name())
		assert.Equal(t, task.KindProducer, task2.Children()[0].Definition().Kind())
		assert.Equal(t, time.Duration(6)*time.Millisecond, task2.Children()[0].Definition().StartAfter())
		assert.Equal(t, time.Duration(200), task2.Children()[0].Definition().Duration())
		assert.Equal(t, 0.4, task2.Children()[0].Definition().FailWithProbability())

		assert.Equal(t, "task2-child2", task2.Children()[1].Definition().Name())
		assert.Equal(t, task.KindProducer, task2.Children()[1].Definition().Kind())
		assert.Equal(t, time.Duration(7)*time.Millisecond, task2.Children()[1].Definition().StartAfter())
		assert.Equal(t, time.Duration(8)*time.Millisecond, task2.Children()[1].Definition().Duration())
		assert.Equal(t, 0.5, task2.Children()[1].Definition().FailWithProbability())

		assert.Equal(t, "task3", result[1].Definition().Name())
		assert.Equal(t, task.KindConsumer, result[1].Definition().Kind())
		assert.Equal(t, time.Duration(100), result[1].Definition().StartAfter())
		assert.Equal(t, time.Duration(200), result[1].Definition().Duration())
		assert.Equal(t, 0.5, result[1].Definition().FailWithProbability())
		assert.Equal(t, "task2-child1", result[1].Definition().LinkedTo()[0].Value())
		assert.Equal(t, "task2-child2", result[1].Definition().LinkedTo()[1].Value())
	})

	t.Run("returns error for invalid task", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				StartAfter: ptrDuration(100),
				Duration:   ptrDuration(200),
				FailWith: FailureCondition{
					Probability: floatPtr(0.5),
				},
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name:       "task1",
							ExternalID: ptrString("^invalid$"),
							Kind:       "client",
						},
					},
				},
			},
		}
		_, err := bp.To()
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid external ID: ^invalid$")
	})
}

func ptrDuration(d time.Duration) *time.Duration {
	return &d
}

func ptrString(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
