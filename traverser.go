package GoHtml

import (
	"sync"
	"github.com/emirpasic/gods/stacks/linkedliststack"
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

//GetCurrentNode returns the current node.
func (t *Traverser) GetCurrentNode() *Node {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	return t.currentNode
}

//SetCurrentNodeTo changes the current node to the newNode.
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
//Walkthrough traverse the node tree from the current node to the end of the node tree by visiting every node. 
func (t *Traverser) Walkthrough(callback func(node *Node)) {
	stack := linkedliststack.New()
	stack.Push(t.GetCurrentNode())

	for stack.Size() > 0{
		currentNode, _ := stack.Pop()
		callback(currentNode.(*Node))
		
		if currentNode.(*Node).GetNextNode() != nil{
			stack.Push(currentNode.(*Node).GetNextNode())
		} 
		if currentNode.(*Node).GetChildNode() != nil{
			stack.Push(currentNode.(*Node).GetChildNode())
		}
	}
}
