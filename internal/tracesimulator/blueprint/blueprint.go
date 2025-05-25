package blueprint

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Blueprint interface defines the method to interpret a blueprint into task trees
type Blueprint interface {
	// Interpret converts the blueprint into a slice of task.TreeNode, which then can be used to create a trace
	Interpret() ([]*task.TreeNode, error)
}
