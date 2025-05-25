package task

import (
	"fmt"
)

// TreeNode represents a node in the task tree
type TreeNode struct {
	definition *Definition
	parent     *TreeNode
	children   []*TreeNode
}

// NewTreeNode creates a new TreeNode with a given ID and Definition
func NewTreeNode(def *Definition) *TreeNode {
	return &TreeNode{
		definition: def,
		children:   make([]*TreeNode, 0),
	}
}

// AddChild adds a child to the current node while ensuring tree constraints
// TODO: Optimize the efficiency of adding a child
// Skipping this for now as it's unlikely to have a large number of children in a single task tree
func (n *TreeNode) AddChild(child *TreeNode) error {
	if child == nil {
		return fmt.Errorf("cannot add nil child")
	}
	if child == n {
		return fmt.Errorf("cannot add self as child")
	}
	if child.parent != nil {
		return fmt.Errorf("child already has a parent: %s", child.parent.definition.name)
	}
	if createsCycle(n, child) {
		return fmt.Errorf("adding %s as child of %s would create a cycle", child.definition.name, n.definition.name)
	}

	child.parent = n
	n.children = append(n.children, child)
	return nil
}

func createsCycle(parent, child *TreeNode) bool {
	// Walk up the chain from parent to root and check if child appears
	for node := parent; node != nil; node = node.parent {
		if node == child {
			return true
		}
	}
	return false
}

func (n *TreeNode) Definition() *Definition {
	return n.definition
}

func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}

func (n *TreeNode) Children() []*TreeNode {
	return n.children
}
