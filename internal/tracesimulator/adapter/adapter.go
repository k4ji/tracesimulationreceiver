package simulator

import "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/span"

// Adapter interface defines the method to transform a tree of spans into a different format
type Adapter[T any] interface {
	// Transform takes the span tree(s) and transforms it into type T
	Transform(rootSpans []*span.TreeNode) (T, error)
}

// NoOpAdapter is a no-operation adapter that does not transform the spans, which can be used for testing or as a placeholder
type NoOpAdapter struct{}

func (a *NoOpAdapter) Transform(rootSpans []*span.TreeNode) ([]*span.TreeNode, error) {
	// No transformation is needed, just return the original spans
	return rootSpans, nil
}
