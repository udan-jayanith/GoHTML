package GoHtml

import (
	"github.com/emirpasic/gods/stacks/linkedliststack"
)

type Traverser struct {
	currentNode *Node
}

// NewTraverser returns a new traverser that can be used to navigate the node tree.
func NewTraverser(startingNode *Node) Traverser {
	return Traverser{
		currentNode: startingNode,
	}
}

// GetCurrentNode returns the current node.
func (t *Traverser) GetCurrentNode() *Node {

	return t.currentNode
}

// SetCurrentNodeTo changes the current node to the newNode.
func (t *Traverser) SetCurrentNodeTo(newNode *Node) {
	t.currentNode = newNode
}

// Next returns the node next to current node and change CurrentNode to the new node.
// Make sure t.currentNode is not nil otherwise program will panic.
func (t *Traverser) Next() *Node {
	currentNode := t.GetCurrentNode()
	t.SetCurrentNodeTo(currentNode.GetNextNode())
	return t.GetCurrentNode()
}

// Previous returns the previous node and change CurrentNode to the new node.
// Make sure t.currentNode is not nil otherwise program will panic.
func (t *Traverser) Previous() *Node {
	currentNode := t.GetCurrentNode()
	t.SetCurrentNodeTo(currentNode.GetPreviousNode())
	return t.GetCurrentNode()
}

type TraverseCondition = bool

const (
	StopWalkthrough     TraverseCondition = false
	ContinueWalkthrough TraverseCondition = true
)

// Walkthrough traverse the node tree from the current node to the end of the node tree by visiting every node.
// Walkthrough traverse the node tree similar to DFS without visiting visited nodes iteratively.
// Walkthrough can be used as a range over iterator or a function that takes a callback and pass every node one by one.
func (t *Traverser) Walkthrough(callback func(node *Node) TraverseCondition) {
	stack := linkedliststack.New()
	if t.GetCurrentNode() == nil {
		return
	}
	stack.Push(t.GetCurrentNode())

	for stack.Size() > 0 {
		currentNode, _ := stack.Pop()
		if !callback(currentNode.(*Node)) {
			return
		}

		if currentNode.(*Node).GetNextNode() != nil {
			stack.Push(currentNode.(*Node).GetNextNode())
		}
		if currentNode.(*Node).GetChildNode() != nil {
			stack.Push(currentNode.(*Node).GetChildNode())
		}
	}
}
