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

// QueryToken store data about basic css selectors(ids, classes, tags).
type QueryToken struct {
	Type         int
	SelectorName string
	Selector     string
}

// TokenizeQuery tokenizes the query and returns a list of QueryToken.
func TokenizeQuery(query string) []QueryToken {
	slice := make([]QueryToken, 0, 1)
	if strings.TrimSpace(query) == "" {
		return slice
	}

	iter := strings.SplitSeq(query, " ")
	for sec := range iter {
		token := QueryToken{}
		switch sec {
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

// matchQueryTokens returns wether the queryTokens match given the node. 
func matchQueryTokens(node *Node, queryTokens []QueryToken) bool {
	if len(queryTokens) == 0 {
		return false
	}
	classList := NewClassList()
	classList.DecodeFrom(node)
	for _, token := range queryTokens {
		switch token.Type {
		case Id:
			idName, _ := node.GetAttribute("id")
			if token.SelectorName != idName {
				return false
			}
		case Tag:
			if node.GetTagName() != token.SelectorName {
				return false
			}
		case Class:
			if !classList.Contains(token.SelectorName) {
				return false
			}
		}
	}
	return true
}

// Query returns the first node that matches with the give query.
func (node *Node) Query(query string) *Node {
	queryTokens := TokenizeQuery(query)

	traverser := NewTraverser(node)
	var res *Node
	traverser.Walkthrough(func(node *Node) TraverseCondition {
		if matchQueryTokens(node, queryTokens) {
			res = node
			return StopWalkthrough
		}
		return ContinueWalkthrough
	})
	return res
}

// QueryAll returns a NodeList containing nodes that matched with the given query.
func (node *Node) QueryAll(query string) NodeList{
	nodeList := NewNodeList()
	queryTokens := TokenizeQuery(query)
	traverser := NewTraverser(node)

	for node := range traverser.Walkthrough{
		if matchQueryTokens(node, queryTokens) {
			nodeList.Append(node)
		}
	}
	return nodeList
}