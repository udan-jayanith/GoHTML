package GoHtml

import (
	"sync"
)

//The DOM Node struct is an abstract base class upon which many other DOM API objects are based, thus letting
// those object types to be used similarly and often interchangeably.
type Node struct {
	nextNode     *Node
	previousNode *Node
	childNode    *Node
	parentNode   *Node

	tagName    string
	attributes map[string]string
	text       string
	rwMutex sync.Mutex
}

//GetNextNode returns node next to the node.
func (node *Node) GetNextNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.nextNode
}

//SetNextNode make nodes next node as nextNode.
func (node *Node) SetNextNode(nextNode *Node){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.nextNode = nextNode
}

//GetPreviousNode returns the previous node.
func (node *Node) GetPreviousNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.previousNode
}

//SetPreviousNode sets nodes previous node to previousNode.
func (node *Node) SetPreviousNode(previousNode *Node){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.previousNode = previousNode
}

//GetChildNode returns the first child elements of this node.
func (node *Node) GetChildNode() *Node{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.childNode
}

//getParentNode returns parent node.
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

//Returns a string with the name of the tag for the given node.
func (node *Node) GetTagName() string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.tagName
}

//SetTagName changes the html tag name to the tagName.
func (node *Node) SetTagName(tagName string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.tagName = tagName
}

//GetAttribute returns the specified attribute form the node.
func (node *Node) GetAttribute(attributeName string) string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.attributes[attributeName]
}

//RemoveAttribute remove or delete the specified attribute.
func (node *Node) RemoveAttribute(attributeName string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	delete(node.attributes, attributeName)
}

//IterateAttributes calls callback at every attribute in the node by passing attribute and value of the each node to the callback.
func (node *Node) IterateAttributes(callback func(attribute, value string)){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	for k, v := range node.attributes{
		callback(k, v)
	}
}

//SetAttribute add a attribute to the node.
func (node *Node) SetAttribute(attribute, value string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.attributes[attribute] = value
}

//GetText returns text on the node. This does not returns text on it's child nodes. If you also wants child nodes text use GetInnerText method on the node.
func (node *Node) GetText() string{
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.text
}

//SetText add text to the node.
func (node *Node) SetText(text string){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.text = text
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

	lastNode := node.GetChildNode().GetLastNode()
	childNode.SetPreviousNode(lastNode)
	lastNode.SetNextNode(childNode)
}

//Append inserts the newNode to end of the node chain.
func (node *Node) Append(newNode *Node) {
	lastNode := node.GetLastNode()
	newNode.SetPreviousNode(lastNode)
	lastNode.SetNextNode(newNode)
}

//GetParent returns a pointer to the parent node.
func (node *Node) GetParent() *Node {
	return node.GetFirstNode().getParentNode()
}

//GetLastNode returns the last node in the node branch.
func (node *Node) GetLastNode() *Node{
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().GetNextNode() != nil {
		traverser.Next()
	}
	return traverser.GetCurrentNode()
}

//GetFirstNode returns the first node of the node branch.
func (node *Node) GetFirstNode() *Node{
	traverser := GetTraverser(node)
	for traverser.GetCurrentNode().GetPreviousNode() != nil{
		traverser.Previous()
	}
	return traverser.GetCurrentNode()
}

//AppendText add text to the node.
func (node *Node) AppendText(text string){
	textNode := CreateNode("")
	textNode.SetText(text)

	if node.GetTagName() == "" || IsVoidTag(node.GetTagName()){
		node.GetLastNode().Append(textNode)
		return
	}
	node.GetLastNode().AppendChild(textNode)
}

//GetInnerText returns all of the text inside the node.
func (node *Node) GetInnerText() string{
	text := ""
	traverser := GetTraverser(node.childNode)
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		if node.GetTagName() != "" {
			return ContinueWalkthrough
		}
		text += node.GetText()
		return ContinueWalkthrough
	})

	return text
}

//RemoveNode removes the node from the branch safely. 
func (node *Node) RemoveNode(){
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	previousNode := node.previousNode
	nextNode := node.nextNode

	node.previousNode = nil
	node.nextNode = nil

	if previousNode != nil {
		previousNode.SetNextNode(nextNode)
	}
	if nextNode != nil {
		nextNode.SetPreviousNode(previousNode)
	}
}

func (node *Node) IsTextNode() bool {
	return node.GetTagName() == ""
}
