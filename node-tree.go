package GoHtml

import (
	"iter"
	"strings"
	"sync"
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
	rwMutex    sync.Mutex
}

// GetNextNode returns node next to the node.
func (node *Node) GetNextNode() *Node {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.nextNode
}

// SetNextNode make nodes next node as nextNode.
func (node *Node) SetNextNode(nextNode *Node) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.nextNode = nextNode
}

// GetPreviousNode returns the previous node.
func (node *Node) GetPreviousNode() *Node {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.previousNode
}

// SetPreviousNode sets nodes previous node to previousNode.
func (node *Node) SetPreviousNode(previousNode *Node) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.previousNode = previousNode
}

// GetChildNode returns the first child elements of this node.
func (node *Node) GetChildNode() *Node {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.childNode
}

// getParentNode returns parent node.
func (node *Node) getParentNode() *Node {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	return node.parentNode
}

func (node *Node) setParentNode(parentNode *Node) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.parentNode = parentNode
}

// Returns a string with the name of the tag for the given node.
func (node *Node) GetTagName() string {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	if strings.ToUpper(node.tagName) == DOCTYPEDTD {
		return strings.ToUpper(node.tagName)
	}

	return node.tagName
}

// SetTagName changes the html tag name to the tagName.
func (node *Node) SetTagName(tagName string) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.tagName = strings.TrimSpace(strings.ToLower(tagName))
}

// GetAttribute returns the specified attribute value form the node. If the specified attribute doesn't exists GetAttribute returns a empty string and false.
func (node *Node) GetAttribute(attributeName string) (string, bool) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	v, ok := node.attributes[attributeName]
	return v, ok
}

// RemoveAttribute remove or delete the specified attribute.
func (node *Node) RemoveAttribute(attributeName string) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	delete(node.attributes, attributeName)
}

// IterAttributes returns a iterator that can be used with range keyword. 
// Where a(first value) is attribute name v(2nd value) is the attribute value.
func (node *Node) IterAttributes() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		node.rwMutex.Lock()
		attributes := node.attributes
		node.rwMutex.Unlock()

		for k, v := range attributes {
			if !yield(k, v) {
				return
			}
		}
	}
}

// SetAttribute add a attribute to the node.
func (node *Node) SetAttribute(attribute, value string) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.attributes[strings.TrimSpace(attribute)] = strings.TrimSpace(value)
}

// GetText returns text on the node. This does not returns text on it's child nodes. If you also wants child nodes text use GetInnerText method on the node.
func (node *Node) GetText() string {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	text := node.text
	text = strings.ReplaceAll(text, "&amp;" ,"&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")

	return text
}

// SetText add text to the node.
func (node *Node) SetText(text string) {
	node.rwMutex.Lock()
	defer node.rwMutex.Unlock()

	node.text = escapeHTML(text)
}

// The AppendChild() method of the Node adds a node to the end of the list of children of a specified parent node.
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
	node.rwMutex.Lock()

	previousNode := node.previousNode
	nextNode := node.nextNode
	parentNode := node.parentNode

	node.previousNode = nil
	node.nextNode = nil
	node.parentNode = nil

	node.rwMutex.Unlock()

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
