package model

import (
	"fmt"
	domainTask "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Service represents a service that executes tasks
type Service struct {
	Name     string
	Resource map[string]string
	Tasks    []Task
}

// To converts the Service to a slice of task.TreeNode
func (s Service) To() ([]*domainTask.TreeNode, error) {
	rootTaskNodes := make([]*domainTask.TreeNode, 0)
	for _, task := range s.Tasks {
		resource := domainTask.NewResource(s.Name, s.Resource)
		rootTaskNode, err := task.ToRootNodeWithResource(resource)
		if err != nil {
			return nil, fmt.Errorf("failed to convert task %s to root node: %w", task.Name, err)
		}
		rootTaskNodes = append(rootTaskNodes, rootTaskNode)
	}
	return rootTaskNodes, nil
}
