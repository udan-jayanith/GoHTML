package GoHtml

import (
	"strings"
)

// GetElementByTagName returns the first node that match with the given tagName by advancing from the node.
func (node *Node) GetElementByTagName(tagName string) *Node {
	tagName = strings.ToLower(strings.TrimSpace(tagName))

	traverser := NewTraverser(node)
	var returnNode *Node
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		if node.GetTagName() == tagName {
			returnNode = node
			return StopWalkthrough
		}

		return ContinueWalkthrough
	})
	return returnNode
}

// GetElementByClassName returns the first node that match with the given className by advancing from the node.
func (node *Node) GetElementByClassName(className string) *Node {
	traverser := NewTraverser(node)
	var returnNode *Node
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		classList := NewClassList()
		classList.DecodeFrom(node)

		if classList.Contains(className) {
			returnNode = node
			return StopWalkthrough
		}
		return ContinueWalkthrough
	})
	return returnNode
}

// GetElementByID returns the first node that match with the given idName by advancing from the node.
func (node *Node) GetElementByID(idName string) *Node {
	traverser := NewTraverser(node)
	var returnNode *Node
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		id, _ := node.GetAttribute("id")
		if id == idName {
			returnNode = node
			return StopWalkthrough
		}
		return ContinueWalkthrough
	})
	return returnNode
}

// GetElementsByClassName returns a NodeList containing nodes that have the given className from the node.
func (node *Node) GetElementsByClassName(className string) NodeList {
	traverser := NewTraverser(node)
	nodeList := NewNodeList()

	traverser.Walkthrough(func(node *Node) TraverseCondition {
		classList := NewClassList()
		classList.DecodeFrom(node)

		if classList.Contains(className) {
			nodeList.Append(node)
		}
		return ContinueWalkthrough
	})
	return nodeList
}

// GetElementsByTagName returns a NodeList containing nodes that have the given tagName from the node.
func (node *Node) GetElementsByTagName(tagName string) NodeList {
	traverser := NewTraverser(node)
	nodeList := NewNodeList()

	traverser.Walkthrough(func(node *Node) TraverseCondition {
		if node.GetTagName() == tagName {
			nodeList.Append(node)
		}
		return ContinueWalkthrough
	})
	return nodeList
}

// GetElementsByClassName returns a NodeList containing nodes that have the given idName from the node.
func (node *Node) GetElementsById(idName string) NodeList {
	traverser := NewTraverser(node)
	nodeList := NewNodeList()

	traverser.Walkthrough(func(node *Node) TraverseCondition {
		id, _ := node.GetAttribute("id")
		if id == idName {
			nodeList.Append(node)
		}
		return ContinueWalkthrough
	})
	return nodeList
}

// Selector types
const (
	Id int = iota
	Tag
	Class
)

// SelectorToken store data about basic css selectors(ids, classes, tags).
type SelectorToken struct {
	Type         int
	SelectorName string
	Selector     string
}

// TokenizeSelector returns a []SelectorToken in selector.
func TokenizeSelector(selector string) []SelectorToken {
	slice := make([]SelectorToken, 0, 1)
	if strings.TrimSpace(selector) == "" {
		return slice
	}

	iter := strings.SplitSeq(selector, " ")
	for sec := range iter {
		token := SelectorToken{}
		switch sec{
		case "", " ", ".", "#":
			continue
		}

		switch string(sec[0]) {
		case ".":
			token.Type = Class
			token.SelectorName = sec[1:]
		case "#":
			token.Type = Id
			token.SelectorName = sec[1:]
		default:
			token.Type = Tag
			token.SelectorName = sec
		}
		token.Selector = sec
		slice = append(slice, token)
	}

	return slice
}

