package GoHtml

import (
	"sync"
)

type Traverser struct {
	currentNode *Node
	rwMutex     sync.Mutex
}

// GetTraverser returns a new traverser that can be used to navigate the node tree.
func GetTraverser(startingNode *Node) Traverser {
	return Traverser{
		currentNode: startingNode,
		rwMutex:     sync.Mutex{},
	}
}

func (t *Traverser) GetCurrentNode() *Node {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	return t.currentNode
}

func (t *Traverser) SetCurrentNodeTo(newNode *Node) {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	t.currentNode = newNode
}

// Next returns the node next to current node and change CurrentNode to the new node.
func (t *Traverser) Next() *Node {
	currentNode := t.GetCurrentNode()
	t.SetCurrentNodeTo(currentNode.GetNextNode())
	return t.GetCurrentNode()
}

// Previous returns the previous node and change CurrentNode to the new node.
func (t *Traverser) Previous() *Node {
	currentNode := t.GetCurrentNode()
	t.SetCurrentNodeTo(currentNode.GetPreviousNode())
	return t.GetCurrentNode()
}

//TODO: use a linked stack
func (t *Traverser) Walkthrough(callback func(node *Node)) {
	stack := []*Node{t.GetCurrentNode()}

	for len(stack) > 0{
		currentNode := stack[len(stack)-1]
		callback(currentNode)

		stack = stack[:len(stack)-1]
		
		if currentNode.GetNextNode() != nil{
			stack = append(stack, currentNode.GetNextNode())
		} 
		if currentNode.GetChildNode() != nil{
			stack = append(stack, currentNode.GetChildNode())
		}
	}
}
