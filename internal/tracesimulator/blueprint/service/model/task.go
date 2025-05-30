package model

import (
	domainTask "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Task represents an operation that can be performed by a service
type Task struct {
	Name                  string
	ExternalID            *domainTask.ExternalID
	Delay                 domainTask.Delay
	Duration              domainTask.Duration
	Kind                  string
	Attributes            map[string]string
	Children              []Task
	ChildOf               *domainTask.ExternalID
	LinkedTo              []*domainTask.ExternalID
	Events                []domainTask.Event
	ConditionalDefinition []domainTask.ConditionalDefinition
}

// ToRootNodeWithResource converts the Task to a root node with the given resource
func (t *Task) ToRootNodeWithResource(resource domainTask.Resource) (*domainTask.TreeNode, error) {
	def := domainTask.NewDefinition(
		t.Name,
		true,
		resource,
		t.Attributes,
		domainTask.FromString(t.Kind),
		t.ExternalID,
		t.Delay,
		t.Duration,
		t.ChildOf,
		t.LinkedTo,
		t.Events,
		t.ConditionalDefinition,
	)
	node := domainTask.NewTreeNode(def)
	for _, child := range t.Children {
		childNode, err := child.toChildNodeWithResource(resource)
		if err != nil {
			return nil, err
		}
		if err = node.AddChild(childNode); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (t *Task) toChildNodeWithResource(resource domainTask.Resource) (*domainTask.TreeNode, error) {
	def := domainTask.NewDefinition(
		t.Name,
		false,
		resource,
		t.Attributes,
		domainTask.FromString(t.Kind),
		t.ExternalID,
		t.Delay,
		t.Duration,
		nil,
		t.LinkedTo,
		t.Events,
		t.ConditionalDefinition,
	)
	node := domainTask.NewTreeNode(def)
	for _, child := range t.Children {
		childNode, err := child.toChildNodeWithResource(resource)
		if err != nil {
			return nil, err
		}
		if err = node.AddChild(childNode); err != nil {
			return nil, err
		}
	}
	return node, nil
}
