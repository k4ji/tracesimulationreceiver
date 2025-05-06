package service

import (
	"github.com/k4ji/tracesimulator/pkg/model/task"
	"github.com/k4ji/tracesimulator/pkg/model/task/taskduration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	t.Run("valid blueprint", func(t *testing.T) {
		duration := Duration{
			Value: ptrString("1000"),
			Mode:  ptrString("absolute"),
		}
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("1000ms"),
					Mode:  ptrString("absolute"),
				},
				Duration: &duration,
				FailWith: FailureCondition{
					Probability: floatPtr(0.5),
				},
			},
			Services: []Service{},
		}
		err := bp.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing delay with default", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("1000ms"),
					Mode:  ptrString("absolute"),
				},
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name: "task1",
							Kind: "client",
							Duration: &Duration{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
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

	t.Run("missing delay without default", func(t *testing.T) {
		bp := &Blueprint{
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name: "task1",
							Kind: "client",
							Duration: &Duration{
								Value: ptrString("1000"),
								Mode:  ptrString("absolute"),
							},
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a delay value or global default")
	})

	t.Run("missing duration with default", func(t *testing.T) {
		duration := Duration{
			Value: ptrString("1000ms"),
			Mode:  ptrString("absolute"),
		}
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &duration,
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name: "task1",
							Delay: &Delay{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "client",
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
							Name: "task1",
							Delay: &Delay{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "client",
							FailWith: FailureCondition{
								Probability: floatPtr(0.5),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a duration value or global default")
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
							Name: "task1",
							Delay: &Delay{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Duration: &Duration{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "client",
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
							Name: "task1",
							Delay: &Delay{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Duration: &Duration{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "client",
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 must have a FailWithProbability value or global default")
	})

	t.Run("invalid delay", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("-1ns"),
					Mode:  ptrString("absolute"),
				},
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name: "task1",
							Kind: "client",
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.ErrorContains(t, err, "invalid delay: absolute delay must be non-negative")
	})

	t.Run("invalid duration", func(t *testing.T) {
		duration := Duration{
			Value: ptrString("-1ms"),
			Mode:  ptrString("absolute"),
		}
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &duration,
			},
			Services: []Service{
				{
					Name: "service1",
					Tasks: []Task{
						{
							Name: "task1",
							Delay: &Delay{
								Value: ptrString("1000ms"),
								Mode:  ptrString("absolute"),
							},
						},
					},
				},
			},
		}
		err := bp.Validate()
		assert.EqualError(t, err, "task task1 has invalid duration: absolute duration must be greater than 0")
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
				Delay: &Delay{
					Value: ptrString("100ns"),
					Mode:  ptrString("absolute"),
				},
				Duration: &Duration{
					Value: ptrString("200ns"),
					Mode:  ptrString("absolute"),
				},
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
							Delay: &Delay{
								Value: ptrString("1ms"),
								Mode:  ptrString("absolute"),
							},
							Duration: &Duration{
								Value: ptrString("2ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "client",
							Attributes: map[string]string{
								"key1": "value1",
							},
							Children: []Task{
								{
									Name: "task1-child",
									Delay: &Delay{
										Value: ptrString("3ms"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Value: ptrString("4ms"),
										Mode:  ptrString("absolute"),
									},
									Kind: "internal",
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
							Name: "task2",
							Duration: &Duration{
								Value: ptrString("5ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "server",
							Children: []Task{
								// missing Value
								{
									Name:       "task2-child1",
									ExternalID: ptrString("task2-child1"),
									Delay: &Delay{
										Value: ptrString("6ms"),
										Mode:  ptrString("absolute"),
									},
									Kind: "producer",
									FailWith: FailureCondition{
										Probability: floatPtr(0.4),
									},
								},
								// missing FailWith
								{
									Name:       "task2-child2",
									ExternalID: ptrString("task2-child2"),
									Delay: &Delay{
										Value: ptrString("7ms"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Value: ptrString("8ms"),
										Mode:  ptrString("absolute"),
									},
									Kind: "producer",
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
						// LinkedTo plus missing StartAfter, Value and FailWith
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
		assert.Equal(t, NewDelayAsAbsoluteDuration(1*time.Millisecond), task1.Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(2)*time.Millisecond), task1.Definition().Duration())
		assert.Equal(t, 0.1, task1.Definition().FailWithProbability())

		assert.Equal(t, "task1-child", task1.Children()[0].Definition().Name())
		assert.Equal(t, "value2", task1.Children()[0].Definition().Attributes()["key2"])
		assert.Equal(t, task.KindInternal, task1.Children()[0].Definition().Kind())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(3)*time.Millisecond), task1.Children()[0].Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(4)*time.Millisecond), task1.Children()[0].Definition().Duration())
		assert.Equal(t, 0.2, task1.Children()[0].Definition().FailWithProbability())

		task2 := result[0].Children()[1]
		assert.Equal(t, "task2", task2.Definition().Name())
		assert.Equal(t, task.KindServer, task2.Definition().Kind())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(100)), task2.Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(5)*time.Millisecond), task2.Definition().Duration())
		assert.Equal(t, "task1", task2.Definition().ChildOf().Value())
		assert.Equal(t, 0.3, task2.Definition().FailWithProbability())

		assert.Equal(t, "task2-child1", task2.Children()[0].Definition().Name())
		assert.Equal(t, task.KindProducer, task2.Children()[0].Definition().Kind())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(6)*time.Millisecond), task2.Children()[0].Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(200)), task2.Children()[0].Definition().Duration())
		assert.Equal(t, 0.4, task2.Children()[0].Definition().FailWithProbability())

		assert.Equal(t, "task2-child2", task2.Children()[1].Definition().Name())
		assert.Equal(t, task.KindProducer, task2.Children()[1].Definition().Kind())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(7)*time.Millisecond), task2.Children()[1].Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Children()[1].Definition().Duration())
		assert.Equal(t, 0.5, task2.Children()[1].Definition().FailWithProbability())

		assert.Equal(t, "task3", result[1].Definition().Name())
		assert.Equal(t, task.KindConsumer, result[1].Definition().Kind())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(100)), result[1].Definition().Delay())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(200)), result[1].Definition().Duration())
		assert.Equal(t, 0.5, result[1].Definition().FailWithProbability())
		assert.Equal(t, "task2-child1", result[1].Definition().LinkedTo()[0].Value())
		assert.Equal(t, "task2-child2", result[1].Definition().LinkedTo()[1].Value())
	})

	t.Run("returns error for invalid task", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("100ns"),
					Mode:  ptrString("absolute"),
				},
				Duration: &Duration{
					Value: ptrString("200ns"),
					Mode:  ptrString("absolute"),
				},
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

func TestConvertConfigWithRelativeAndAbsoluteDelayModes(t *testing.T) {
	t.Run("without default delay", func(t *testing.T) {
		{
			bp := &Blueprint{
				Default: DefaultValues{
					Duration: &Duration{
						Value: ptrString("10ms"),
						Mode:  ptrString("absolute"),
					},
					FailWith: FailureCondition{
						Probability: floatPtr(0.0),
					},
				},
				Services: []Service{
					{
						Name: "service",
						Tasks: []Task{
							{
								Name: "root",
								Delay: &Delay{
									Value: ptrString("10ms"),
									Mode:  ptrString("absolute"),
								},
								Kind: "internal",
								Children: []Task{
									{
										Name: "task-with-explicit-relative-duration-and-mode",
										Delay: &Delay{
											Value: ptrString("0.9"),
											Mode:  ptrString("relative"),
										},
										Kind: "internal",
										Children: []Task{
											{
												Name: "task-with-explicit-absolute-duration-and-mode",
												Delay: &Delay{
													Value: ptrString("8ms"),
													Mode:  ptrString("absolute"),
												},
												Kind: "internal",
											},
										},
									},
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
			assert.Len(t, result, 1)
			// check the first task tree
			root := result[0]
			assert.Equal(t, "root", root.Definition().Name())
			assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Delay())

			// check the first child task
			task1 := root.Children()[0]
			assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
			assert.Equal(t, NewDelayAsRelativeDuration(0.9), task1.Definition().Delay())

			// check the second child task
			task2 := task1.Children()[0]
			assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
			assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Delay())
		}
	})

	t.Run("with default delay of relative mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("0.1"),
					Mode:  ptrString("relative"),
				},
				Duration: &Duration{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Delay: &Delay{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Delay: &Delay{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Delay: &Delay{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-relative-duration",
													Delay: &Delay{
														Value: ptrString("0.7"),
													},
													Kind: "internal",
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Delay())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDelayAsRelativeDuration(0.9), task1.Definition().Delay())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Delay())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration", task3.Definition().Name())
		assert.Equal(t, NewDelayAsRelativeDuration(0.7), task3.Definition().Delay())
	})

	t.Run("with default delay of absolute mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("1ms"),
					Mode:  ptrString("absolute"),
				},
				Duration: &Duration{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Delay: &Delay{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Delay: &Delay{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Delay: &Delay{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-duration",
													Delay: &Delay{
														Value: ptrString("7ms"),
													},
													Kind: "internal",
													Children: []Task{
														{
															Name: "task-with-explicit-absolute-mode",
															Delay: &Delay{
																Mode: ptrString("absolute"),
															},
															Kind: "internal",
															Children: []Task{
																{
																	Name: "task-without-explicit-duration",
																	Kind: "internal",
																},
															},
														},
													},
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Delay())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDelayAsRelativeDuration(0.9), task1.Definition().Delay())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Delay())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration", task3.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(7)*time.Millisecond), task3.Definition().Delay())

		// check the child of the third child task
		task4 := task3.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-mode", task4.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task4.Definition().Delay())

		// check the child of the fourth child task
		task5 := task4.Children()[0]
		assert.Equal(t, "task-without-explicit-duration", task5.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task5.Definition().Delay())
	})

	t.Run("with default delay of duration", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Value: ptrString("1ms"),
				},
				Duration: &Duration{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Delay: &Delay{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Delay: &Delay{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Delay: &Delay{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-mode",
													Delay: &Delay{
														Mode: ptrString("absolute"),
													},
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Delay())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDelayAsRelativeDuration(0.9), task1.Definition().Delay())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Delay())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-mode", task3.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task3.Definition().Delay())
	})

	t.Run("with default delay of mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Delay: &Delay{
					Mode: ptrString("absolute"),
				},
				Duration: &Duration{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Delay: &Delay{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Delay: &Delay{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Delay: &Delay{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-duration",
													Delay: &Delay{
														Value: ptrString("7ms"),
													},
													Kind: "internal",
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Delay())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDelayAsRelativeDuration(0.9), task1.Definition().Delay())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Delay())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration", task3.Definition().Name())
		assert.Equal(t, NewDelayAsAbsoluteDuration(time.Duration(7)*time.Millisecond), task3.Definition().Delay())
	})

	t.Run("invalid combinations of delay duration and mode", func(t *testing.T) {
		testCases := []struct {
			name      string
			blueprint Blueprint
			errorMsg  string
		}{
			{
				name: "duration (in absolute) and relative mode",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task1",
									Delay: &Delay{
										Value: ptrString("10ms"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "invalid relative delay format",
			},
			{
				name: "duration (in relative) and absolute mode",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task2",
									Delay: &Delay{
										Value: ptrString("0.5"),
										Mode:  ptrString("absolute"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "invalid absolute delay format",
			},
			{
				name: "duration only",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task3",
									Delay: &Delay{
										Value: ptrString("10ms"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "missing required field: delay.mode",
			},
			{
				name: "mode only",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task4",
									Delay: &Delay{
										Mode: ptrString("absolute"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "missing required field: delay.for",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := tc.blueprint.To()
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorMsg)
			})
		}
	})
}

func TestConvertConfigWithRelativeAndAbsoluteDurationModes(t *testing.T) {
	t.Run("without default duration", func(t *testing.T) {
		{
			bp := &Blueprint{
				Default: DefaultValues{
					Delay: &Delay{
						Value: ptrString("10ms"),
						Mode:  ptrString("absolute"),
					},
					FailWith: FailureCondition{
						Probability: floatPtr(0.0),
					},
				},
				Services: []Service{
					{
						Name: "service",
						Tasks: []Task{
							{
								Name: "root",
								Duration: &Duration{
									Value: ptrString("10ms"),
									Mode:  ptrString("absolute"),
								},
								Kind: "internal",
								Children: []Task{
									{
										Name: "task-with-explicit-relative-duration-and-mode",
										Duration: &Duration{
											Value: ptrString("0.9"),
											Mode:  ptrString("relative"),
										},
										Kind: "internal",
										Children: []Task{
											{
												Name: "task-with-explicit-absolute-duration-and-mode",
												Duration: &Duration{
													Value: ptrString("8ms"),
													Mode:  ptrString("absolute"),
												},
												Kind: "internal",
											},
										},
									},
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
			assert.Len(t, result, 1)
			// check the first task tree
			root := result[0]
			assert.Equal(t, "root", root.Definition().Name())
			assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Duration())

			// check the first child task
			task1 := root.Children()[0]
			assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
			assert.Equal(t, NewDurationAsRelativeDuration(0.9), task1.Definition().Duration())

			// check the second child task
			task2 := task1.Children()[0]
			assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
			assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Duration())
		}
	})

	t.Run("with default duration of relative mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &Duration{
					Value: ptrString("0.1"),
					Mode:  ptrString("relative"),
				},
				Delay: &Delay{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Duration: &Duration{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Duration: &Duration{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Duration: &Duration{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-relative-duration",
													Duration: &Duration{
														Value: ptrString("0.7"),
													},
													Kind: "internal",
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Duration())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDurationAsRelativeDuration(0.9), task1.Definition().Duration())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Duration())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration", task3.Definition().Name())
		assert.Equal(t, NewDurationAsRelativeDuration(0.7), task3.Definition().Duration())
	})

	t.Run("with default duration of absolute mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &Duration{
					Value: ptrString("1ms"),
					Mode:  ptrString("absolute"),
				},
				Delay: &Delay{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Duration: &Duration{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Duration: &Duration{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Duration: &Duration{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-duration",
													Duration: &Duration{
														Value: ptrString("7ms"),
													},
													Kind: "internal",
													Children: []Task{
														{
															Name: "task-with-explicit-absolute-mode",
															Duration: &Duration{
																Mode: ptrString("absolute"),
															},
															Kind: "internal",
															Children: []Task{
																{
																	Name: "task-without-explicit-duration",
																	Kind: "internal",
																},
															},
														},
													},
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Duration())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDurationAsRelativeDuration(0.9), task1.Definition().Duration())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Duration())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration", task3.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(7)*time.Millisecond), task3.Definition().Duration())

		// check the child of the third child task
		task4 := task3.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-mode", task4.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task4.Definition().Duration())

		// check the child of the fourth child task
		task5 := task4.Children()[0]
		assert.Equal(t, "task-without-explicit-duration", task5.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task5.Definition().Duration())
	})

	t.Run("with default duration of duration", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &Duration{
					Value: ptrString("1ms"),
				},
				Delay: &Delay{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Duration: &Duration{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Duration: &Duration{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Duration: &Duration{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-mode",
													Duration: &Duration{
														Mode: ptrString("absolute"),
													},
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Duration())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDurationAsRelativeDuration(0.9), task1.Definition().Duration())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Duration())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-mode", task3.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(1)*time.Millisecond), task3.Definition().Duration())
	})

	t.Run("with default duration of mode", func(t *testing.T) {
		bp := &Blueprint{
			Default: DefaultValues{
				Duration: &Duration{
					Mode: ptrString("absolute"),
				},
				Delay: &Delay{
					Value: ptrString("10ms"),
					Mode:  ptrString("absolute"),
				},
				FailWith: FailureCondition{
					Probability: floatPtr(0.0),
				},
			},
			Services: []Service{
				{
					Name: "service",
					Tasks: []Task{
						{
							Name: "root",
							Duration: &Duration{
								Value: ptrString("10ms"),
								Mode:  ptrString("absolute"),
							},
							Kind: "internal",
							Children: []Task{
								{
									Name: "task-with-explicit-relative-duration-and-mode",
									Duration: &Duration{
										Value: ptrString("0.9"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
									Children: []Task{
										{
											Name: "task-with-explicit-absolute-duration-and-mode",
											Duration: &Duration{
												Value: ptrString("8ms"),
												Mode:  ptrString("absolute"),
											},
											Kind: "internal",
											Children: []Task{
												{
													Name: "task-with-explicit-absolute-duration",
													Duration: &Duration{
														Value: ptrString("7ms"),
													},
													Kind: "internal",
												},
											},
										},
									},
								},
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
		assert.Len(t, result, 1)
		// check the first task tree
		root := result[0]
		assert.Equal(t, "root", root.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(10)*time.Millisecond), root.Definition().Duration())

		// check the first child task
		task1 := root.Children()[0]
		assert.Equal(t, "task-with-explicit-relative-duration-and-mode", task1.Definition().Name())
		assert.Equal(t, NewDurationAsRelativeDuration(0.9), task1.Definition().Duration())

		// check the second child task
		task2 := task1.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration-and-mode", task2.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(8)*time.Millisecond), task2.Definition().Duration())

		// check the child of the second child task
		task3 := task2.Children()[0]
		assert.Equal(t, "task-with-explicit-absolute-duration", task3.Definition().Name())
		assert.Equal(t, NewDurationAsAbsoluteDuration(time.Duration(7)*time.Millisecond), task3.Definition().Duration())
	})

	t.Run("invalid combinations of duration value and mode", func(t *testing.T) {
		testCases := []struct {
			name      string
			blueprint Blueprint
			errorMsg  string
		}{
			{
				name: "duration (in absolute) and relative mode",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task1",
									Delay: &Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Value: ptrString("10ms"),
										Mode:  ptrString("relative"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "invalid relative duration format",
			},
			{
				name: "duration (in relative) and absolute mode",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task2",
									Delay: &Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Value: ptrString("0.5"),
										Mode:  ptrString("absolute"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "invalid absolute duration format",
			},
			{
				name: "duration only",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task3",
									Delay: &Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Value: ptrString("10ms"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "missing required field: duration.mode",
			},
			{
				name: "mode only",
				blueprint: Blueprint{
					Services: []Service{
						{
							Name: "service",
							Tasks: []Task{
								{
									Name: "task4",
									Delay: &Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &Duration{
										Mode: ptrString("absolute"),
									},
									Kind: "internal",
								},
							},
						},
					},
				},
				errorMsg: "missing required field: duration.for",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := tc.blueprint.To()
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorMsg)
			})
		}
	})
}

func ptrString(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func NewDelayAsAbsoluteDuration(td time.Duration) task.Delay {
	ad, _ := taskduration.NewAbsoluteDuration(td)
	d, _ := task.NewDelay(ad)
	return *d
}

func NewDelayAsRelativeDuration(v float64) task.Delay {
	rd, _ := taskduration.NewRelativeDuration(v)
	d, _ := task.NewDelay(rd)
	return *d
}

func NewDurationAsAbsoluteDuration(td time.Duration) task.Duration {
	ad, _ := taskduration.NewAbsoluteDuration(td)
	d, _ := task.NewDuration(ad)
	return *d
}

func NewDurationAsRelativeDuration(v float64) task.Duration {
	rd, _ := taskduration.NewRelativeDuration(v)
	d, _ := task.NewDuration(rd)
	return *d
}
