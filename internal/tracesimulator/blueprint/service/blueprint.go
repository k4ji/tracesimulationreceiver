package service

import (
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint/service/model"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Blueprint implements `Blueprint` interface
var _ blueprint.Blueprint = (*Blueprint)(nil)

// Blueprint represents a blueprint based on tasks grouped by services
type Blueprint struct {
	services []model.Service
}

// NewServiceBlueprint creates a new service blueprint
func NewServiceBlueprint(services []model.Service) Blueprint {
	return Blueprint{
		services: services,
	}
}

func (sb *Blueprint) Interpret() ([]*task.TreeNode, error) {
	rootTaskNodes := make([]*task.TreeNode, 0)
	TasksByExternalID := make(map[task.ExternalID]*task.TreeNode)

	// Convert each service to trees of tasks
	for _, service := range sb.services {
		serviceRootTaskNodes, err := service.To()
		if err != nil {
			return nil, fmt.Errorf("failed to convert service %s to task tree: %w", service.Name, err)
		}
		rootTaskNodes = append(rootTaskNodes, serviceRootTaskNodes...)
	}

	// Add all task nodes to the map by ExternalID while checking for duplicates
	var addNodesToMap func(node *task.TreeNode) error
	addNodesToMap = func(node *task.TreeNode) error {
		if node.Definition().ExternalID() != nil {
			if _, exists := TasksByExternalID[*node.Definition().ExternalID()]; exists {
				return fmt.Errorf("duplicate ExternalID detected: %s", *node.Definition().ExternalID())
			}
			TasksByExternalID[*node.Definition().ExternalID()] = node
		}
		for _, child := range node.Children() {
			if err := addNodesToMap(child); err != nil {
				return err
			}
		}
		return nil
	}
	for _, rootTaskNode := range rootTaskNodes {
		if err := addNodesToMap(rootTaskNode); err != nil {
			return nil, err
		}
	}

	// Connect task nodes across different services according to their parent-child relationships
	traceRootTaskNodes := make([]*task.TreeNode, 0)
	for _, rootTaskNode := range rootTaskNodes {
		if rootTaskNode.Definition().ChildOf() != nil {
			parentSpan := TasksByExternalID[*rootTaskNode.Definition().ChildOf()]
			if parentSpan == nil {
				return nil, fmt.Errorf("parent task not found for %s", *rootTaskNode.Definition().ChildOf())
			}
			if err := parentSpan.AddChild(rootTaskNode); err != nil {
				return nil, fmt.Errorf("failed to add child %s to parent %s: %w", rootTaskNode.Definition().Name(), parentSpan.Definition().Name(), err)
			}
		} else {
			traceRootTaskNodes = append(traceRootTaskNodes, rootTaskNode)
		}
	}

	return traceRootTaskNodes, nil
}
