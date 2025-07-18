package GoHtml

import (
	"sync"
)

type Node struct {
	NextNode     *Node
	PreviousNode *Node
	ChildNode    *Node
	parentNode   *Node

	TagName    string
	Attributes map[string]string
	AttributesMutex sync.Mutex
	Closed     bool
	Text       string
	RWMutex sync.Mutex
}

//The AppendChild() method of the Node adds a node to the end of the list of children of a specified parent node.
func (node *Node) AppendChild(childNode *Node) {
	if node.ChildNode == nil {
		node.ChildNode = childNode
		childNode.parentNode = node
		return
	}

	traverser := GetTraverser(node.ChildNode)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	traverser.CurrentNode.NextNode = childNode
	childNode.PreviousNode = traverser.CurrentNode
}

//Append inserts the newNode to end of the node chain.
func (node *Node) Append(newNode *Node) {
	traverser := GetTraverser(node)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	newNode.PreviousNode = traverser.CurrentNode
	traverser.CurrentNode.NextNode = newNode
}

//GetParent returns a pointer to the parent node.
func (node *Node) GetParent() *Node {
	traverser := GetTraverser(node)
	for traverser.CurrentNode.parentNode == nil && traverser.CurrentNode.PreviousNode != nil {
		traverser.Previous()
	}
	return traverser.CurrentNode.parentNode
}

//GetLastNode returns the last node in the node chain.
func (node *Node) GetLastNode() *Node{
	traverser := GetTraverser(node)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	return traverser.CurrentNode
}

//GetFirstNode returns the first node of the node chain.
func (node *Node) GetFirstNode() *Node{
	traverser := GetTraverser(node)
	for traverser.CurrentNode.PreviousNode != nil{
		traverser.Previous()
	}
	return traverser.CurrentNode
}

//AppendText add text to the node
func (node *Node) AppendText(text string){
	textNode := CreateNode("")
	textNode.Text = text

	node.GetLastNode().NextNode = textNode
}

//GetInnerText returns all of the text in the node excluding child nodes text.
func (node *Node) GetInnerText() string{
	text := ""
	traverser := GetTraverser(node.ChildNode)
	for traverser.CurrentNode != nil{
		text += traverser.CurrentNode.Text
		traverser.Next()
	}

	return text
}
