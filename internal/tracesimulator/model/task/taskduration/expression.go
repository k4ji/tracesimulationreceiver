package taskduration

import "time"

// Expression is an interface that represents a duration expression.
type Expression interface {
	// Resolve resolves the expression to a time.Duration based on the provided context.
	Resolve(context interface{}) (*time.Duration, error)
}
