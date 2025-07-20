package GoHtml

import (
	"strings"
	"sync"
)

//CreateNode returns a initialized new node.
func CreateNode(tagName string) *Node {
	return &Node{
		tagName: strings.ToLower(strings.TrimSpace(tagName)),
		rwMutex: sync.Mutex{},
		closed:  true,
	}
}

//DeepCloneNode clones the node without having references to it's original parent node, previous node and next node.
func DeepCloneNode(node *Node) *Node{
	newNode := Node{
		childNode: node.GetChildNode(),
		tagName: node.GetTagName(),
		attributes: node.GetAttributes(),
		text: node.GetText(),
		closed: node.isClosed(),

		rwMutex: sync.Mutex{},
	}

	return &newNode
}

//CloneNode copy the node.
func CloneNode(node *Node) *Node{ 
	newNode := DeepCloneNode(node)
	newNode.setParentNode(node.getParentNode())
	newNode.SetPreviousNode(node.GetPreviousNode())
	newNode.SetNextNode(node.GetNextNode())

	return newNode
}

//ApplySaveChanges replaces the nodes previous and parent node with the given node.
func ApplySaveChanges(node *Node){
	previousNode := node.GetPreviousNode()
	if previousNode != nil{
		previousNode.SetNextNode(node)
	}

	parentNode := node.getParentNode()
	if parentNode != nil{
		parentNode.rwMutex.Lock()
		parentNode.childNode = node
		parentNode.rwMutex.Unlock()
	}
}