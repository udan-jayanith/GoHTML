package GoHtml

import (
	"strings"
	"golang.org/x/net/html"
)

// Node is a struct that represents a html elements. Nodes can have sibling nodes(NextNode and Previous Node) and child node that represent the child elements.
// Text is also stored as a node which can be checked by using IsTextNode method.
type Node struct {
	nextNode     *Node
	previousNode *Node
	childNode    *Node
	parentNode   *Node

	tagName    string
	attributes map[string]string
	text       string
}

// GetNextNode returns node next to the node.
func (node *Node) GetNextNode() *Node {
	return node.nextNode
}

// SetNextNode make nodes next node as nextNode.
func (node *Node) SetNextNode(nextNode *Node) {
	node.nextNode = nextNode
}

// GetPreviousNode returns the previous node.
func (node *Node) GetPreviousNode() *Node {
	return node.previousNode
}

// SetPreviousNode sets nodes previous node to previousNode.
func (node *Node) SetPreviousNode(previousNode *Node) {
	node.previousNode = previousNode
}

// GetChildNode returns the first child elements of this node.
func (node *Node) GetChildNode() *Node {
	return node.childNode
}

// getParentNode returns parent node.
func (node *Node) getParentNode() *Node {
	return node.parentNode
}

func (node *Node) setParentNode(parentNode *Node) {
	node.parentNode = parentNode
}

// Returns a string with the name of the tag for the given node.
func (node *Node) GetTagName() string {
	if strings.ToUpper(node.tagName) == DOCTYPEDTD {
		return strings.ToUpper(node.tagName)
	}

	return node.tagName
}

// SetTagName changes the html tag name to the tagName.
func (node *Node) SetTagName(tagName string) {
	node.tagName = strings.TrimSpace(strings.ToLower(tagName))
}

// GetAttribute returns the specified attribute value form the node. If the specified attribute doesn't exists GetAttribute returns a empty string and false.
func (node *Node) GetAttribute(attributeName string) (string, bool) {
	v, ok := node.attributes[strings.TrimSpace(strings.ToLower(attributeName))]
	return v, ok
}

// RemoveAttribute remove or delete the specified attribute.
func (node *Node) RemoveAttribute(attributeName string) {
	delete(node.attributes, strings.TrimSpace(strings.ToLower(attributeName)))
	
}

// IterateAttributes calls callback at every attribute in the node by passing attribute and value of the node.
func (node *Node) IterateAttributes(callback func(attribute, value string)) {
	attributes := node.attributes
	for k, v := range attributes {
		callback(k, v)
	}
}

// SetAttribute add a attribute to the node.
func (node *Node) SetAttribute(attribute, value string) {
	node.attributes[strings.ToLower(strings.TrimSpace(attribute))] = strings.TrimSpace(value)
}

// GetText returns text on the node. This does not returns text on it's child nodes. If you also wants child nodes text use GetInnerText method on the node.
// HTML tags in returns value get escaped.
func (node *Node) GetText() string {
	text := node.text
	return text
}

// SetText add text to the node.
// SetText unescapes entities like "&lt;" to become "<".
func (node *Node) SetText(text string) {
	node.text = html.UnescapeString(text)
}

// The AppendChild() method of the Node adds a node to the end of the list of children of a specified parent node.
func (node *Node) AppendChild(childNode *Node) {
	if node.GetChildNode() == nil {
		node.childNode = childNode
		childNode.setParentNode(node)
		return
	}

	lastNode := node.GetChildNode().GetLastNode()
	childNode.SetPreviousNode(lastNode)
	lastNode.SetNextNode(childNode)
}

// Append inserts the newNode to end of the node chain.
func (node *Node) Append(newNode *Node) {
	lastNode := node.GetLastNode()
	newNode.SetPreviousNode(lastNode)
	lastNode.SetNextNode(newNode)
}

// GetParent returns a pointer to the parent node.
func (node *Node) GetParent() *Node {
	return node.GetFirstNode().getParentNode()
}

// GetLastNode returns the last node in the node branch.
func (node *Node) GetLastNode() *Node {
	traverser := NewTraverser(node)
	for traverser.GetCurrentNode().GetNextNode() != nil {
		traverser.Next()
	}
	return traverser.GetCurrentNode()
}

// GetFirstNode returns the first node of the node branch.
func (node *Node) GetFirstNode() *Node {
	traverser := NewTraverser(node)
	for traverser.GetCurrentNode().GetPreviousNode() != nil {
		traverser.Previous()
	}
	return traverser.GetCurrentNode()
}

// AppendText append text to the node.
func (node *Node) AppendText(text string) {
	textNode := CreateTextNode(text)
	if node.GetTagName() == "" || IsVoidTag(node.GetTagName()) {
		node.GetLastNode().Append(textNode)
		return
	}
	node.GetLastNode().AppendChild(textNode)
}

// GetInnerText returns all of the text inside the node.
func (node *Node) GetInnerText() string {
	text := ""
	traverser := NewTraverser(node.childNode)
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		if node.GetTagName() != "" {
			return ContinueWalkthrough
		}
		text += node.GetText()
		return ContinueWalkthrough
	})

	return text
}

// RemoveNode removes the node from the branch safely by connecting sibling nodes.
func (node *Node) RemoveNode() {
	previousNode := node.previousNode
	nextNode := node.nextNode
	parentNode := node.parentNode

	node.previousNode = nil
	node.nextNode = nil
	node.parentNode = nil

	if previousNode != nil {
		previousNode.SetNextNode(nextNode)
	}
	if nextNode != nil {
		nextNode.SetPreviousNode(previousNode)
	}

	if nextNode != nil && previousNode == nil {
		nextNode.setParentNode(parentNode)
	}

	if parentNode != nil {
		parentNode.childNode = nextNode
	}

}

// IsTextNode returns a boolean value indicating node is a text node or not.
func (node *Node) IsTextNode() bool {
	return node.GetTagName() == ""
}
