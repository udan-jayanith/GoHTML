package GoHtml

import (
	"sync"
)

type Node struct {
	nextNode     *Node
	previousNode *Node
	childNode    *Node
	parentNode   *Node

	tagName    string
	attributes map[string]string
	closed     bool
	text       string
	rwMutex sync.Mutex
}

func (node *Node) GetNextNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.nextNode
}

func (node *Node) SetNextNode(nextNode *Node){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.nextNode = nextNode
}

func (node *Node) GetPreviousNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.previousNode
}

func (node *Node) SetPreviousNode(previousNode *Node){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.previousNode = previousNode
}

func (node *Node) GetChildNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.childNode
}

func (node *Node) getParentNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.parentNode
}

func (node *Node) setParentNode(parentNode *Node){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.parentNode = parentNode
}

func (node *Node) GetTagName() string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.tagName
}

func (node *Node) SetTagName(tagName string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.tagName = tagName
}

func (node *Node) GetAttribute(attributeName string) string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.attributes[attributeName]
}

func (node *Node) GetAttributes() map[string]string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.attributes
}

func (node *Node) SetAttribute(attribute, value string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.attributes[attribute] = value
}

func (node *Node) GetText() string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.text
}

func (node *Node) SetText(text string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.text = text
}

func (node *Node) isClosed() bool{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.closed
}

func (node *Node) setClosedTo(closed bool){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.closed = closed
}

//The AppendChild() method of the Node adds a node to the end of the list of children of a specified parent node.
func (node *Node) AppendChild(childNode *Node) {
	if node.GetChildNode() == nil {
		node.rwMutex.Lock()
		node.childNode = childNode
		node.rwMutex.Unlock()

		childNode.setParentNode(node)
		return
	}

	traverser := GetTraverser(node.GetChildNode())
	for traverser.GetCurrentNode().GetNextNode() != nil {
		traverser.Next()
	}
	traverser.GetCurrentNode().SetNextNode(childNode)
	childNode.SetPreviousNode(traverser.GetCurrentNode())
}

//Append inserts the newNode to end of the node chain.
func (node *Node) Append(newNode *Node) {
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().GetNextNode() != nil {
		traverser.Next()
	}
	newNode.SetPreviousNode(traverser.GetCurrentNode())
	traverser.GetCurrentNode().SetNextNode(newNode)
}

//GetParent returns a pointer to the parent node.
func (node *Node) GetParent() *Node {
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().getParentNode() == nil && traverser.GetCurrentNode().GetPreviousNode() != nil {
		traverser.Previous()
	}
	return traverser.GetCurrentNode().getParentNode()
}

//GetLastNode returns the last node in the node chain.
func (node *Node) GetLastNode() *Node{
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().GetNextNode() != nil {
		traverser.Next()
	}
	return traverser.GetCurrentNode()
}

//GetFirstNode returns the first node of the node chain.
func (node *Node) GetFirstNode() *Node{
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().GetPreviousNode() != nil{
		traverser.Previous()
	}
	return traverser.GetCurrentNode()
}

//AppendText add text to the node
func (node *Node) AppendText(text string){
	textNode := CreateNode("")
	textNode.SetText(text)

	node.GetLastNode().AppendChild(textNode)
}

//GetInnerText returns all of the text in the node excluding child nodes text.
func (node *Node) GetInnerText() string{
	text := ""
	traverser := GetTraverser(node.childNode)
	for traverser.GetCurrentNode() != nil{
		text += traverser.GetCurrentNode().GetText()
		traverser.Next()
	}

	return text
}
