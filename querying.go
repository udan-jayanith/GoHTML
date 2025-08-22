package GoHtml

import (
	//"iter"
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
QuerySearch tokenizes the query string and search for nodes that matches with the right most query token. After matching right most query it proceeds to match nodes parents nodes for left over tokens and then passed that node to (yield/range). QuerySearch search the whole node tree for matches unless yield get canceled or range iterator get cancel.
*/
/*
func QuerySearch(node *Node, query string) iter.Seq[*Node] {
	traverser := NewTraverser(node)
	return func(yield func(node *Node) bool) {
		queryTokens := TokenizeQuery(query)
		iter := traverser.Walkthrough
		for node := range iter {
			i := matchFromRightMostQueryToken(node, queryTokens, len(queryTokens)-1)
			if i == len(queryTokens)-1{
				continue
			}
			parentNode := node.GetParent()
			for parentNode != nil && i>=0 {
				i = matchFromRightMostQueryToken(parentNode, queryTokens, i)
				parentNode = parentNode.GetParent()
			}
			if i < 0 && !yield(node){
				return
			}
		}

	}
}

// matchFromRightMostQueryToken tries to match query tokens from right to left and return the index at which point query token last matched.
func matchFromRightMostQueryToken(node *Node, queryTokens []QueryToken, i int) int {
	classList := NewClassList()
	classList.DecodeFrom(node)
	checked := make(map[string]struct{})
outer:
	for i >= 0 {
		token := queryTokens[i]
		_, ok := checked[token.Selector]
		if ok {
			break
		} else {
			checked[token.Selector] = struct{}{}
		}

		switch token.Type {
		case Id:
			idName, _ := node.GetAttribute("id")
			if token.SelectorName != idName {
				break outer
			}
		case Class:
			if !classList.Contains(token.SelectorName) {
				break outer
			}
		case Tag:
			if node.GetTagName() != token.SelectorName {
				break outer
			}
		}
		i--
	}
	return i
}

// QuerySelector only returns the first node that matches with the QuerySearch.
func (node *Node) QuerySelector(query string) *Node {
	iter := QuerySearch(node, query)
	for node := range iter {
		return node
	}
	return nil
}

// QuerySelectorAll stores nodes passed down by QuerySearch in a nodeList and returns the nodeList.
func (node *Node) QuerySelectorAll(query string) NodeList {
	iter := QuerySearch(node, query)
	nodeList := NewNodeList()

	for node := range iter {
		nodeList.Append(node)
	}
	return nodeList
}

*/