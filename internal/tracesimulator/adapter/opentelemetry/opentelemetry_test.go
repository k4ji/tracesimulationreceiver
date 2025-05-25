package opentelemetry

import (
	simulator "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.27.0"
	mathRand "math/rand"
	"testing"
	"time"
)

// TODO create each node instance directly
//
//	for now, this test uses simulator to create the nodes and get the otel spans
//	since it's not possible to create span.Node, which is located in other package
func TestAdapter_Transform(t *testing.T) {
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
				"resource-key-service-a": "resource-value-service-a",
			},
			Tasks: []model.Task{
				{
					Name:       "root-task-a",
					ExternalID: rootTaskAExternalID,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(1000 * time.Millisecond),
					Kind:       "client",
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
						"attribute-key-root-task-a": "attribute-value-root-task-a",
					},
					Children: []model.Task{
						{
							Name:       "child-task-a1",
							ExternalID: childTaskA1ExternalID,
							Delay:      NewAbsoluteDurationDelay(500 * time.Millisecond),
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
							Attributes: map[string]string{
								"attribute-key-child-task-a1": "attribute-value-child-task-a1",
							},
						},
						{
							Name:       "child-task-a2",
							ExternalID: childTaskA2ExternalID,
							Delay:      NewAbsoluteDurationDelay(1000 * time.Millisecond),
							Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
							Kind:       "client",
							Events: []task.Event{
								task.NewEvent(
									"event-child-task-a2-1",
									NewAbsoluteDurationDelay(0),
									map[string]string{
										"attribute-key-event-child-task-a2-1": "attribute-value-event-child-task-a2-1",
									},
								),
							},
							Attributes: map[string]string{
								"attribute-key-child-task-a2": "attribute-value-child-task-a2",
							},
						},
					},
				},
			},
		},
		// Linked spans
		{
			Name: "service-b",
			Resource: map[string]string{
				"resource-key-service-b": "resource-value-service-b",
			},
			Tasks: []model.Task{
				{
					Name:       "root-task-b",
					ExternalID: rootTaskBExternalID,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(2000 * time.Millisecond),
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
					Attributes: map[string]string{
						"attribute-key-root-task-b": "attribute-value-root-task-b",
					},
					Children: []model.Task{
						{
							Name:       "child-task-b1",
							ExternalID: nil,
							Delay:      NewAbsoluteDurationDelay(1000 * time.Millisecond),
							Duration:   NewAbsoluteDurationDuration(500 * time.Millisecond),
							Kind:       "internal",
							Events: []task.Event{
								task.NewEvent(
									"event-child-task-b1-1",
									NewAbsoluteDurationDelay(0),
									map[string]string{
										"attribute-key-event-child-task-b1-1": "attribute-value-event-child-task-b1-1",
									},
								),
							},
							Attributes: map[string]string{
								"attribute-key-child-task-b1": "attribute-value-child-task-b1",
							},
						},
					},
					LinkedTo: []*task.ExternalID{
						childTaskA1ExternalID,
						rootTaskCExternalID,
					},
				},
			},
		},
		// child span of another service span
		{
			Name: "service-c",
			Resource: map[string]string{
				"resource-key-service-c": "resource-value-service-c",
			},
			Tasks: []model.Task{
				{
					Name:       "root-task-c",
					ExternalID: rootTaskCExternalID,
					Delay:      NewAbsoluteDurationDelay(1000 * time.Millisecond),
					Duration:   NewAbsoluteDurationDuration(2000 * time.Millisecond),
					Kind:       "server",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-c-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-c-1": "attribute-value-event-root-task-c-1",
							},
						),
					},
					ChildOf: rootTaskAExternalID,
					Attributes: map[string]string{
						"attribute-key-root-task-c": "attribute-value-root-task-c",
					},
				},
			},
		},
		// error spans
		{
			Name: "service-d",
			Resource: map[string]string{
				"resource-key-service-d": "resource-value-service-d",
			},
			Tasks: []model.Task{
				{
					Name:       "root-task-d",
					ExternalID: nil,
					Delay:      NewAbsoluteDurationDelay(0),
					Duration:   NewAbsoluteDurationDuration(4000 * time.Millisecond),
					Kind:       "server",
					Events: []task.Event{
						task.NewEvent(
							"event-root-task-d-1",
							NewAbsoluteDurationDelay(0),
							map[string]string{
								"attribute-key-event-root-task-d-1": "attribute-value-event-root-task-d-1",
							},
						),
					},
					ChildOf: childTaskA2ExternalID,
					Attributes: map[string]string{
						"attribute-key-root-task-d": "attribute-value-root-task-d",
					},
					ConditionalDefinition: []*task.ConditionalDefinition{
						task.NewConditionalDefinition(
							task.NewProbabilisticCondition(1.0, mathRand.Float64),
							[]task.Effect{
								task.FromMarkAsFailedEffect(task.NewMarkAsFailedEffect(ptrString("error"))),
							},
						),
					},
				},
			},
		},
	})

	sim := simulator.New[[]ptrace.Traces](NewAdapter())
	traces, err := sim.Run(&blueprint, now)

	// Create a map of span names to their corresponding span for easy lookup
	spanMap := make(map[string]ptrace.Span)
	for _, trace := range traces {
		rs := trace.ResourceSpans()
		for i := 0; i < rs.Len(); i++ {
			scopeSpans := rs.At(i).ScopeSpans()
			for j := 0; j < scopeSpans.Len(); j++ {
				spans := scopeSpans.At(j).Spans()
				for k := 0; k < spans.Len(); k++ {
					span := spans.At(k)
					spanMap[span.Name()] = span
				}
			}
		}
	}

	t.Run("transform spans to otel spans", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Len(t, traces, 2)

		assert.Equal(t, traces[0].SpanCount(), 5)
		assert.Equal(t, traces[1].SpanCount(), 2)
	})

	t.Run("resource attributes are kept", func(t *testing.T) {
		expectedResources := map[string]map[string]string{
			"service-a": {
				"resource-key-service-a": "resource-value-service-a",
			},
			"service-b": {
				"resource-key-service-b": "resource-value-service-b",
			},
			"service-c": {
				"resource-key-service-c": "resource-value-service-c",
			},
			"service-d": {
				"resource-key-service-d": "resource-value-service-d",
			},
		}

		for _, trace := range traces {
			rs := trace.ResourceSpans()
			for i := 0; i < rs.Len(); i++ {
				resource := rs.At(i).Resource()
				serviceName, _ := resource.Attributes().Get(conventions.AttributeServiceName)
				expectedResource, exists := expectedResources[serviceName.AsString()]
				assert.True(t, exists, "Unexpected service name: %s", serviceName.AsString())
				for key, value := range expectedResource {
					attr, ok := resource.Attributes().Get(key)
					assert.True(t, ok, "Missing resource attribute: %s", key)
					assert.Equal(t, value, attr.AsString())
				}
			}
		}
	})

	t.Run("span names, kinds, and attributes are kept", func(t *testing.T) {
		expectedSpans := map[string]struct {
			Name       string
			Kind       ptrace.SpanKind
			Attributes map[string]string
		}{
			"root-task-a": {
				Name: "root-task-a",
				Kind: ptrace.SpanKindClient,
				Attributes: map[string]string{
					"attribute-key-root-task-a": "attribute-value-root-task-a",
				},
			},
			"child-task-a1": {
				Name: "child-task-a1",
				Kind: ptrace.SpanKindProducer,
				Attributes: map[string]string{
					"attribute-key-child-task-a1": "attribute-value-child-task-a1",
				},
			},
			"child-task-a2": {
				Name: "child-task-a2",
				Kind: ptrace.SpanKindClient,
				Attributes: map[string]string{
					"attribute-key-child-task-a2": "attribute-value-child-task-a2",
				},
			},
			"root-task-b": {
				Name: "root-task-b",
				Kind: ptrace.SpanKindConsumer,
				Attributes: map[string]string{
					"attribute-key-root-task-b": "attribute-value-root-task-b",
				},
			},
			"child-task-b1": {
				Name: "child-task-b1",
				Kind: ptrace.SpanKindInternal,
				Attributes: map[string]string{
					"attribute-key-child-task-b1": "attribute-value-child-task-b1",
				},
			},
			"root-task-c": {
				Name: "root-task-c",
				Kind: ptrace.SpanKindServer,
				Attributes: map[string]string{
					"attribute-key-root-task-c": "attribute-value-root-task-c",
				},
			},
			"root-task-d": {
				Name: "root-task-d",
				Kind: ptrace.SpanKindServer,
				Attributes: map[string]string{
					"attribute-key-root-task-d": "attribute-value-root-task-d",
				},
			},
		}

		for _, trace := range traces {
			rs := trace.ResourceSpans()
			for i := 0; i < rs.Len(); i++ {
				scopeSpans := rs.At(i).ScopeSpans()
				for j := 0; j < scopeSpans.Len(); j++ {
					spans := scopeSpans.At(j).Spans()
					for k := 0; k < spans.Len(); k++ {
						span := spans.At(k)
						expected, exists := expectedSpans[span.Name()]
						assert.True(t, exists, "Unexpected span: %s", span.Name())
						assert.Equal(t, expected.Name, span.Name())
						assert.Equal(t, expected.Kind, span.Kind())
						for key, value := range expected.Attributes {
							attr, ok := span.Attributes().Get(key)
							assert.True(t, ok, "Missing attribute: %s", key)
							assert.Equal(t, value, attr.AsString())
						}
					}
				}
			}
		}
	})

	t.Run("span start and end times are kept", func(t *testing.T) {
		expectedStartTime := now.Add(-5000 * time.Millisecond)

		expectedTimes := map[string]struct {
			Start time.Time
			End   time.Time
		}{
			"root-task-a": {
				Start: expectedStartTime,
				End:   expectedStartTime.Add(1000 * time.Millisecond),
			},
			"child-task-a1": {
				Start: expectedStartTime.Add(500 * time.Millisecond),
				End:   expectedStartTime.Add(1000 * time.Millisecond),
			},
			"child-task-a2": {
				Start: expectedStartTime.Add(1000 * time.Millisecond),
				End:   expectedStartTime.Add(1500 * time.Millisecond),
			},
			"root-task-b": {
				Start: expectedStartTime,
				End:   expectedStartTime.Add(2000 * time.Millisecond),
			},
			"child-task-b1": {
				Start: expectedStartTime.Add(1000 * time.Millisecond),
				End:   expectedStartTime.Add(1500 * time.Millisecond),
			},
			"root-task-c": {
				Start: expectedStartTime.Add(1000 * time.Millisecond),
				End:   expectedStartTime.Add(3000 * time.Millisecond),
			},
			"root-task-d": {
				Start: expectedStartTime.Add(1000 * time.Millisecond),
				End:   expectedStartTime.Add(5000 * time.Millisecond),
			},
		}

		for spanName, expectedTime := range expectedTimes {
			span, spanExists := spanMap[spanName]
			assert.True(t, spanExists, "Span not found: %s", spanName)

			if spanExists {
				actualStart := span.StartTimestamp().AsTime()
				actualEnd := span.EndTimestamp().AsTime()

				assert.Equal(t, expectedTime.Start.In(time.UTC), actualStart, "Unexpected end time for span: %s", spanName)
				assert.Equal(t, expectedTime.End.In(time.UTC), actualEnd, "Unexpected end time for span: %s", spanName)
			}
		}
	})

	t.Run("span events are kept", func(t *testing.T) {
		expectedStartTime := now.Add(-5000 * time.Millisecond)

		expectedEvents := map[string][]struct {
			Name       string
			OccurredAt time.Time
			Attributes map[string]string
		}{
			"root-task-a": {
				{
					Name:       "event-root-task-a-1",
					OccurredAt: expectedStartTime,
					Attributes: map[string]string{
						"attribute-key-event-root-task-a-1": "attribute-value-event-root-task-a-1",
					},
				},
				{
					Name:       "event-root-task-a-2",
					OccurredAt: expectedStartTime.Add(100 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-root-task-a-2": "attribute-value-event-root-task-a-2",
					},
				},
			},
			"child-task-a1": {
				{
					Name:       "event-child-task-a1-1",
					OccurredAt: expectedStartTime.Add(500 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-child-task-a1-1": "attribute-value-event-child-task-a1-1",
					},
				},
			},
			"child-task-a2": {
				{
					Name:       "event-child-task-a2-1",
					OccurredAt: expectedStartTime.Add(1000 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-child-task-a2-1": "attribute-value-event-child-task-a2-1",
					},
				},
			},
			"root-task-b": {
				{
					Name:       "event-root-task-b-1",
					OccurredAt: expectedStartTime,
					Attributes: map[string]string{
						"attribute-key-event-root-task-b-1": "attribute-value-event-root-task-b-1",
					},
				},
			},
			"child-task-b1": {
				{
					Name:       "event-child-task-b1-1",
					OccurredAt: expectedStartTime.Add(1000 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-child-task-b1-1": "attribute-value-event-child-task-b1-1",
					},
				},
			},
			"root-task-c": {
				{
					Name:       "event-root-task-c-1",
					OccurredAt: expectedStartTime.Add(1000 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-root-task-c-1": "attribute-value-event-root-task-c-1",
					},
				},
			},
			"root-task-d": {
				{
					Name:       "event-root-task-d-1",
					OccurredAt: expectedStartTime.Add(1000 * time.Millisecond),
					Attributes: map[string]string{
						"attribute-key-event-root-task-d-1": "attribute-value-event-root-task-d-1",
					},
				},
			},
		}
		for spanName, expectedEvents := range expectedEvents {
			span, spanExists := spanMap[spanName]
			assert.True(t, spanExists, "Span not found: %s", spanName)

			if spanExists {
				events := span.Events()
				assert.Equal(t, len(expectedEvents), events.Len(), "Unexpected number of events for span: %s", spanName)

				for i := 0; i < events.Len(); i++ {
					event := events.At(i)
					expectedEvent := expectedEvents[i]

					assert.Equal(t, expectedEvent.Name, event.Name())
					assert.Equal(t, expectedEvent.OccurredAt.In(time.UTC), event.Timestamp().AsTime(), "Unexpected event time for span: %s", spanName)
					for key, value := range expectedEvent.Attributes {
						attr, ok := event.Attributes().Get(key)
						assert.True(t, ok, "Missing attribute: %s", key)
						assert.Equal(t, value, attr.AsString())
					}
				}
			}
		}
	})

	t.Run("parent-child relationships are kept", func(t *testing.T) {
		expectedChildParent := map[string]string{
			"child-task-a1": "root-task-a",
			"child-task-a2": "root-task-a",
			"child-task-b1": "root-task-b",
			"root-task-c":   "root-task-a",
			"root-task-d":   "child-task-a2",
		}

		for child, parent := range expectedChildParent {
			childSpan, childExists := spanMap[child]
			parentSpan, parentExists := spanMap[parent]

			assert.True(t, childExists, "Child span not found: %s", child)
			assert.True(t, parentExists, "Parent span not found: %s", parent)

			if childExists && parentExists {
				assert.Equal(t, parentSpan.SpanID(), childSpan.ParentSpanID(), "Parent span ID mismatch for child: %s", child)
				assert.Equal(t, parentSpan.TraceID(), childSpan.TraceID(), "Parent and child should have the same TraceID for child: %s", child)
			}
		}
	})

	t.Run("span links are kept", func(t *testing.T) {
		expectedLinks := map[string][]string{
			"root-task-b": {"child-task-a1", "root-task-c"},
		}

		for spanName, linkedSpanNames := range expectedLinks {
			span, spanExists := spanMap[spanName]
			assert.True(t, spanExists, "Span not found: %s", spanName)

			if spanExists {
				actualLinks := make(map[string]struct{})
				links := span.Links()
				for i := 0; i < links.Len(); i++ {
					link := links.At(i)
					for linkedName, linkedSpan := range spanMap {
						if link.TraceID() == linkedSpan.TraceID() && link.SpanID() == linkedSpan.SpanID() {
							actualLinks[linkedName] = struct{}{}
						}
					}
				}

				for _, linkedName := range linkedSpanNames {
					_, exists := actualLinks[linkedName]
					assert.True(t, exists, "Expected link not found: %s -> %s", spanName, linkedName)
				}
			}
		}
	})

	t.Run("status code is kept", func(t *testing.T) {
		expectedStatuses := map[string]ptrace.StatusCode{
			"root-task-a":   ptrace.StatusCodeUnset,
			"child-task-a1": ptrace.StatusCodeUnset,
			"child-task-a2": ptrace.StatusCodeUnset,
			"root-task-b":   ptrace.StatusCodeUnset,
			"child-task-b1": ptrace.StatusCodeUnset,
			"root-task-c":   ptrace.StatusCodeUnset,
			"root-task-d":   ptrace.StatusCodeError,
		}

		for spanName, expectedStatus := range expectedStatuses {
			span, spanExists := spanMap[spanName]
			assert.True(t, spanExists, "Span not found: %s", spanName)

			if spanExists {
				actualStatus := span.Status().Code()
				assert.Equal(t, expectedStatus, actualStatus, "Unexpected status for span: %s", spanName)
			}
		}
	})

	t.Run("status message is kept", func(t *testing.T) {
		expectedStatuses := map[string]string{
			"root-task-a":   "",
			"child-task-a1": "",
			"child-task-a2": "",
			"root-task-b":   "",
			"child-task-b1": "",
			"root-task-c":   "",
			"root-task-d":   "error",
		}

		for spanName, expectedStatus := range expectedStatuses {
			span, spanExists := spanMap[spanName]
			assert.True(t, spanExists, "Span not found: %s", spanName)

			if spanExists {
				actualStatus := span.Status().Message()
				assert.Equal(t, expectedStatus, actualStatus, "Unexpected status message for span: %s", spanName)
			}
		}
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
