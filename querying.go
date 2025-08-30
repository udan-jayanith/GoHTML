package GoHtml

import (
	"iter"
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

/*
QuerySearch search returns a iterator that traverse through the node tree from given node and passes nodes that matches the given selector.
*/
func QuerySearch(node *Node, selector string) iter.Seq[*Node] {
	traverser := NewTraverser(node)
	return func(yield func(node *Node) bool) {
		selectorTokens := TokenizeSelectorsAndCombinators(selector)
		for node := range traverser.Walkthrough {
			if matchFromRightMostSelectors(node, selectorTokens) && !yield(node) {
				return
			}
		}

	}
}

// matchFromRightMostQueryToken tries to match query tokens from right to left and return the index at which point query token last matched.
func matchFromRightMostSelectors(node *Node, selectorTokens []CombinatorEl) bool {
	for i := len(selectorTokens) - 1; i >= 0; i-- {
		if node == nil {
			break
		}
		node = selectorTokens[i].getMatchingNode(node)
	}
	return node != nil
}

// QuerySelector returns the first node that matches with the selector from the node.
func (node *Node) QuerySelector(selector string) *Node {
	iter := QuerySearch(node, selector)
	for node := range iter {
		return node
	}
	return nil
}

// QuerySelectorAll returns a NodeList that has node that matches the selector form the node.
func (node *Node) QuerySelectorAll(selector string) NodeList {
	iter := QuerySearch(node, selector)
	nodeList := NewNodeList()

	for node := range iter {
		nodeList.Append(node)
	}
	return nodeList
}
