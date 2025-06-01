package span

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"testing"
	"time"

	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"github.com/stretchr/testify/assert"
)

func TestFromTaskTree(t *testing.T) {
	type testCase struct {
		name        string
		taskTree    *task.TreeNode
		traceID     TraceID
		baseEndTime time.Time
		idGen       func() ID
		randGen     func() float64
		expected    *TreeNode
	}

	baseTime := time.Now()
	traceID := NewTraceID([16]byte{0x01})

	testCases := []testCase{
		{
			name: "transform a root task to a span",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
						map[string]string{"team": "team-a"},
						task.KindServer,
						func() *task.ExternalID { id, _ := task.NewExternalID("root-task"); return id }(),
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{
							task.NewEvent(
								"root-task-event",
								NewAbsoluteDurationDelay(1*time.Second),
								make(map[string]string),
							),
						},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(1.0, func() float64 { return 1.0 }),
								[]task.Effect{
									task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
								},
							),
						})
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				resource:             task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
				attributes:           map[string]string{"team": "team-a"},
				kind:                 KindServer,
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				parentID:             nil,
				externalID:           func() *task.ExternalID { id, _ := task.NewExternalID("root-task"); return id }(),
				status:               StatusOK,
				children:             []*TreeNode{},
				linkedTo:             []*TreeNode{},
				events: []Event{
					NewEvent("root-task-event", baseTime.Add(2*time.Second), make(map[string]string)),
				},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "copy resource and attributes of a task into a span",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
							map[string]string{"key1": "val1"},
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(1*time.Second),
							NewAbsoluteDurationDuration(2*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
								map[string]string{"key2": "val2"},
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(3*time.Second),
								NewAbsoluteDurationDuration(4*time.Second),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
				attributes:           map[string]string{"key1": "val1"},
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				events:               []Event{},
				status:               StatusOK,
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", map[string]string{"service.version": "1.0.0"}),
						attributes:           map[string]string{"key2": "val2"},
						startTime:            baseTime.Add(4 * time.Second),
						endTime:              baseTime.Add(8 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "set start time relative to the parent span's start time",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(1*time.Second),
							NewAbsoluteDurationDuration(2*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", make(map[string]string)),
								make(map[string]string),
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(3*time.Second),
								NewAbsoluteDurationDuration(4*time.Second),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			randGen: func() float64 { return 0.0 },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				events:               []Event{},
				status:               StatusOK,
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           make(map[string]string),
						startTime:            baseTime.Add(4 * time.Second),
						endTime:              baseTime.Add(8 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "set delay relative to the parent span's duration",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(0),
							NewAbsoluteDurationDuration(10*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", make(map[string]string)),
								make(map[string]string),
								task.KindClient,
								nil,
								NewRelativeDurationDelay(0.5),
								NewAbsoluteDurationDuration(20*time.Second),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			randGen: func() float64 { return 0.0 },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime,
				endTime:              baseTime.Add(10 * time.Second),
				events:               []Event{},
				status:               StatusOK,
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           make(map[string]string),
						startTime:            baseTime.Add(5 * time.Second),
						endTime:              baseTime.Add(25 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "set event delay relative to the parent span's duration",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(0),
							NewAbsoluteDurationDuration(30*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{
								task.NewEvent(
									"relative-delay-event",
									NewRelativeDurationDelay(0.5),
									make(map[string]string),
								),
							},
							[]task.ConditionalDefinition{},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", make(map[string]string)),
								make(map[string]string),
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(20*time.Second),
								NewAbsoluteDurationDuration(10*time.Second),
								nil,
								[]*task.ExternalID{},
								[]task.Event{
									task.NewEvent(
										"absolute-delay-event",
										NewAbsoluteDurationDelay(5*time.Second),
										make(map[string]string),
									),
								},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			randGen: func() float64 { return 0.0 },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime,
				endTime:              baseTime.Add(30 * time.Second),
				events: []Event{
					NewEvent("relative-delay-event", baseTime.Add(15*time.Second), make(map[string]string)),
				},
				status: StatusOK,
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           make(map[string]string),
						startTime:            baseTime.Add(20 * time.Second),
						endTime:              baseTime.Add(30 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events: []Event{
							NewEvent("absolute-delay-event", baseTime.Add(25*time.Second), make(map[string]string)),
						},
						linkedToExternalID: []*task.ExternalID{},
						children:           []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "set duration relative to the parent span's duration",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(0),
							NewAbsoluteDurationDuration(10*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", make(map[string]string)),
								make(map[string]string),
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(0),
								NewRelativeDurationDuration(0.5),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime,
				endTime:              baseTime.Add(10 * time.Second),
				events:               []Event{},
				status:               StatusOK,
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           make(map[string]string),
						startTime:            baseTime.Add(0 * time.Second),
						endTime:              baseTime.Add(5 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "record events based on probability",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(0.5, func() float64 { return 0.4 }),
								[]task.Effect{
									task.FromRecordEventEffect(task.NewRecordEventEffect(
										task.NewEvent(
											"event-name",
											NewAbsoluteDurationDelay(1*time.Second),
											map[string]string{"key": "value"},
										),
									)),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusOK,
				linkedTo:             []*TreeNode{},
				events: []Event{
					NewEvent("event-name", baseTime.Add(2*time.Second), map[string]string{"key": "value"}),
				},
				linkedToExternalID: []*task.ExternalID{},
				children:           []*TreeNode{},
			},
		},
		{
			name: "not record events based on probability",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(0.5, func() float64 { return 0.6 }),
								[]task.Effect{
									task.FromRecordEventEffect(task.NewRecordEventEffect(
										task.NewEvent(
											"event-name",
											NewAbsoluteDurationDelay(1*time.Second),
											map[string]string{"key": "value"},
										),
									)),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusOK,
				linkedTo:             []*TreeNode{},
				events:               []Event{},
				linkedToExternalID:   []*task.ExternalID{},
				children:             []*TreeNode{},
			},
		},
		{
			name: "generate error spans",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(0.5, func() float64 { return 0.4 }),
								[]task.Effect{
									task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusError("error"),
				linkedTo:             []*TreeNode{},
				events:               []Event{},
				linkedToExternalID:   []*task.ExternalID{},
				children:             []*TreeNode{},
			},
		},
		{
			name: "annotate span with attributes",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						map[string]string{"key1": "val1"},
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(1.0, func() float64 { return 0.0 }),
								[]task.Effect{
									task.FromAnnotateEffect(task.NewAnnotateEffect(
										map[string]string{"key2": "val2"},
									)),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           map[string]string{"key1": "val1", "key2": "val2"},
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusOK,
				linkedTo:             []*TreeNode{},
				events:               []Event{},
				linkedToExternalID:   []*task.ExternalID{},
				children:             []*TreeNode{},
			},
		},
		{
			name: "annotate span without attributes",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						nil,
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(1.0, func() float64 { return 0.0 }),
								[]task.Effect{
									task.FromAnnotateEffect(task.NewAnnotateEffect(
										map[string]string{"key": "value"},
									)),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           map[string]string{"key": "value"},
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusOK,
				linkedTo:             []*TreeNode{},
				events:               []Event{},
				linkedToExternalID:   []*task.ExternalID{},
				children:             []*TreeNode{},
			},
		},
		{
			name: "evaluate nested condition",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(0),
							NewAbsoluteDurationDuration(10*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{
								task.NewConditionalDefinition(
									task.NewAtLeastCondition(
										1,
										task.NewChildCondition(task.NewHasAttributeCondition("key1")),
									),
									[]task.Effect{
										task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
									},
								),
							},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task-1",
								false,
								task.NewResource("service-a", make(map[string]string)),
								map[string]string{"key1": "val1"},
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(0),
								NewRelativeDurationDuration(0.5),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task-2",
								false,
								task.NewResource("service-a", make(map[string]string)),
								map[string]string{"key2": "val2"},
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(0),
								NewRelativeDurationDuration(0.5),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child-1 span
					{0x03}, // ID for the child-2 span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime,
				endTime:              baseTime.Add(10 * time.Second),
				events:               []Event{},
				status:               StatusError("error"),
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task-1",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           map[string]string{"key1": "val1"},
						startTime:            baseTime.Add(0 * time.Second),
						endTime:              baseTime.Add(5 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
					{
						id:                   NewSpanID([8]byte{0x03}),
						traceID:              traceID,
						name:                 "child-task-2",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           map[string]string{"key2": "val2"},
						startTime:            baseTime.Add(0 * time.Second),
						endTime:              baseTime.Add(5 * time.Second),
						status:               StatusOK,
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "conditions are evaluated from below to top",
			taskTree: func() *task.TreeNode {
				root := task.NewTreeNode(
					func() task.Definition {
						def := task.NewDefinition(
							"root-task",
							true,
							task.NewResource("service-a", make(map[string]string)),
							make(map[string]string),
							task.KindInternal,
							nil,
							NewAbsoluteDurationDelay(0),
							NewAbsoluteDurationDuration(10*time.Second),
							nil,
							[]*task.ExternalID{},
							[]task.Event{},
							[]task.ConditionalDefinition{
								task.NewConditionalDefinition(
									task.NewAtLeastCondition(
										1,
										task.NewChildCondition(task.NewMarkedAsFailedCondition()),
									),
									[]task.Effect{
										task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
									},
								),
							},
						)
						return def
					}(),
				)
				//nolint:errcheck
				root.AddChild(
					task.NewTreeNode(
						func() task.Definition {
							def := task.NewDefinition(
								"child-task",
								false,
								task.NewResource("service-a", make(map[string]string)),
								make(map[string]string),
								task.KindClient,
								nil,
								NewAbsoluteDurationDelay(0),
								NewRelativeDurationDuration(0.5),
								nil,
								[]*task.ExternalID{},
								[]task.Event{},
								[]task.ConditionalDefinition{
									task.NewConditionalDefinition(
										task.NewProbabilisticCondition(1.0, func() float64 { return 0.0 }),
										[]task.Effect{
											task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
										},
									),
								},
							)
							return def
						}(),
					),
				)
				return root
			}(),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen: func() func() ID {
				ids := [][8]byte{
					{0x01}, // ID for the root span
					{0x02}, // ID for the child-1 span
				}
				index := 0
				return func() ID {
					id := NewSpanID(ids[index])
					index++
					return id
				}
			}(),
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime,
				endTime:              baseTime.Add(10 * time.Second),
				events:               []Event{},
				status:               StatusError("error"),
				children: []*TreeNode{
					{
						id:                   NewSpanID([8]byte{0x02}),
						traceID:              traceID,
						name:                 "child-task",
						isResourceEntryPoint: false,
						kind:                 KindClient,
						resource:             task.NewResource("service-a", make(map[string]string)),
						attributes:           make(map[string]string),
						startTime:            baseTime.Add(0 * time.Second),
						endTime:              baseTime.Add(5 * time.Second),
						status:               StatusError("error"),
						parentID:             func() *ID { id := NewSpanID([8]byte{0x01}); return &id }(),
						externalID:           nil,
						linkedTo:             []*TreeNode{},
						events:               []Event{},
						linkedToExternalID:   []*task.ExternalID{},
						children:             []*TreeNode{},
					},
				},
				linkedTo:           []*TreeNode{},
				linkedToExternalID: []*task.ExternalID{},
			},
		},
		{
			name: "apply multiple effects",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(1*time.Second),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{
							task.NewConditionalDefinition(
								task.NewProbabilisticCondition(1.0, func() float64 { return 0.0 }),
								[]task.Effect{
									task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect("error")),
									task.FromRecordEventEffect(task.NewRecordEventEffect(
										task.NewEvent(
											"event-name",
											NewAbsoluteDurationDelay(1*time.Second),
											map[string]string{"key": "value"},
										),
									)),
								},
							),
						},
					)
					return def
				}(),
			),
			traceID:     traceID,
			baseEndTime: baseTime,
			idGen:       func() ID { return NewSpanID([8]byte{0x01}) },
			expected: &TreeNode{
				id:                   NewSpanID([8]byte{0x01}),
				traceID:              traceID,
				name:                 "root-task",
				isResourceEntryPoint: true,
				kind:                 KindInternal,
				resource:             task.NewResource("service-a", make(map[string]string)),
				attributes:           make(map[string]string),
				startTime:            baseTime.Add(1 * time.Second),
				endTime:              baseTime.Add(3 * time.Second),
				status:               StatusError("error"),
				linkedTo:             []*TreeNode{},
				events: []Event{
					NewEvent("event-name", baseTime.Add(2*time.Second), map[string]string{"key": "value"}),
				},
				linkedToExternalID: []*task.ExternalID{},
				children:           []*TreeNode{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			span, err := FromTaskTree(tc.taskTree, tc.traceID, tc.baseEndTime, tc.idGen)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, span)
		})
	}
}

func TestFromTaskTreeError(t *testing.T) {
	type testCase struct {
		name     string
		taskTree *task.TreeNode
		traceID  TraceID
	}

	testCases := []testCase{
		{
			name: "error when relative duration is specified but no parent task is provided",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewRelativeDurationDelay(0.5),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{},
						[]task.ConditionalDefinition{},
					)
					return def
				}(),
			),
			traceID: NewTraceID([16]byte{0x01}),
		},
		{
			name: "error when event delay is greater than the task duration",
			taskTree: task.NewTreeNode(
				func() task.Definition {
					def := task.NewDefinition(
						"root-task",
						true,
						task.NewResource("service-a", make(map[string]string)),
						make(map[string]string),
						task.KindInternal,
						nil,
						NewAbsoluteDurationDelay(0),
						NewAbsoluteDurationDuration(2*time.Second),
						nil,
						[]*task.ExternalID{},
						[]task.Event{
							task.NewEvent(
								"event-name",
								NewAbsoluteDurationDelay(3*time.Second),
								make(map[string]string),
							),
						},
						[]task.ConditionalDefinition{},
					)
					return def
				}(),
			),
			traceID: NewTraceID([16]byte{0x01}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FromTaskTree(tc.taskTree, tc.traceID, time.Now(), func() ID { return NewSpanID([8]byte{0x01}) })
			assert.Error(t, err)
		})
	}
}

func TestShiftTimestamps(t *testing.T) {
	now := time.Now()
	rootNodeStartTime := now.Add(0 * time.Second)
	rootNodeEndTime := now.Add(1 * time.Second)
	childNodeStartTime := now.Add(2 * time.Second)
	childNodeEndTime := now.Add(3 * time.Second)
	node := &TreeNode{
		startTime: rootNodeStartTime,
		endTime:   rootNodeEndTime,
		children: []*TreeNode{
			{
				startTime: childNodeStartTime,
				endTime:   childNodeEndTime,
			},
		},
	}

	delta := 3 * time.Second
	node.ShiftTimestamps(delta)

	assert.Equal(t, rootNodeStartTime.Add(delta), node.startTime)
	assert.Equal(t, rootNodeEndTime.Add(delta), node.endTime)
	assert.Equal(t, childNodeStartTime.Add(delta), node.children[0].startTime)
	assert.Equal(t, childNodeEndTime.Add(delta), node.children[0].endTime)
}

func NewAbsoluteDurationDelay(duration time.Duration) task.Delay {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := task.NewDelay(e)
	return *d
}

func NewRelativeDurationDelay(v float64) task.Delay {
	e, _ := taskduration.NewRelativeDuration(v)
	d, _ := task.NewDelay(e)
	return *d
}

func NewAbsoluteDurationDuration(duration time.Duration) task.Duration {
	e, _ := taskduration.NewAbsoluteDuration(duration)
	d, _ := task.NewDuration(e)
	return *d
}

func NewRelativeDurationDuration(v float64) task.Duration {
	e, _ := taskduration.NewRelativeDuration(v)
	d, _ := task.NewDuration(e)
	return *d
}
