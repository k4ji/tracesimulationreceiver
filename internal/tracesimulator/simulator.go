package simulator

import (
	"crypto/rand"
	"fmt"
	simulator "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/adapter"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/span"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"time"
)

// Simulator is a struct that simulates traces based on a blueprint and export them to a specific format using an adapter.
type Simulator[T any] struct {
	adapter simulator.Adapter[T]
}

// New creates a new Simulator instance with the provided adapter.
func New[T any](adapter simulator.Adapter[T]) *Simulator[T] {
	return &Simulator[T]{adapter: adapter}
}

// Run executes the simulation by interpreting the blueprint, generating spans, and transforming them using the adapter.
func (s *Simulator[T]) Run(blueprint blueprint.Blueprint, baseEndTime time.Time) (T, error) {
	var zero T
	traceRootTaskNodes, err := blueprint.Interpret()
	if err != nil {
		return zero, fmt.Errorf("failed to interpret blueprint: %w", err)
	}

	// Convert task trees to spans and hold mapping of ExternalID to span
	rootSpans := make([]*span.TreeNode, 0, len(traceRootTaskNodes))
	externalIDToSpan := make(map[task.ExternalID]*span.TreeNode)
	for _, taskTree := range traceRootTaskNodes {
		traceID := generateTraceID()
		rootSpan, err := span.FromTaskTree(taskTree, traceID, baseEndTime, generateSpanID)
		if err != nil {
			return zero, fmt.Errorf("failed to construct span tree: %w", err)
		}
		mp := rootSpan.ExternalIDToSpan()
		for externalID, spanNode := range mp {
			if _, exists := externalIDToSpan[externalID]; exists {
				return zero, fmt.Errorf("failed to construct span tree: duplicate ExternalID detected, {%s}", externalID)
			}
			externalIDToSpan[externalID] = spanNode
		}
		rootSpans = append(rootSpans, rootSpan)
	}

	// Link spans to their parents based on ExternalID
	// This must be done after all spans are created since the linked spans may not be created yet
	for _, rootSpan := range rootSpans {
		err := rootSpan.LinkSpan(externalIDToSpan)
		if err != nil {
			return zero, fmt.Errorf("failed to link spans: %w", err)
		}
	}

	// Shift timestamps to ensure all spans end before the current time
	latestEndTime := baseEndTime
	for _, rootSpan := range rootSpans {
		latestEndTime = s.findLatestEndTime(rootSpan, latestEndTime)
	}
	adjustmentDuration := baseEndTime.Sub(latestEndTime)
	for _, rootSpan := range rootSpans {
		rootSpan.ShiftTimestamps(adjustmentDuration)
	}

	// Convert spans to the desired format using the adapter
	transformed, err := s.adapter.Transform(rootSpans)
	if err != nil {
		return zero, fmt.Errorf("failed to transform spans: %w", err)
	}

	return transformed, nil
}

func (s *Simulator[T]) findLatestEndTime(node *span.TreeNode, latestEndTime time.Time) time.Time {
	if node.EndTime().After(latestEndTime) {
		latestEndTime = node.EndTime()
	}
	for _, child := range node.Children() {
		latestEndTime = s.findLatestEndTime(child, latestEndTime)
	}
	return latestEndTime
}

func generateTraceID() span.TraceID {
	var id [16]byte
	_, _ = rand.Read(id[:])
	return span.NewTraceID(id)
}

func generateSpanID() span.ID {
	var id [8]byte
	_, _ = rand.Read(id[:])
	return span.NewSpanID(id)
}
