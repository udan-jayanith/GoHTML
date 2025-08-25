package GoHtml

import (
	"strings"

	"golang.org/x/net/html"
)

type BasicSelector int

const (
	Id BasicSelector = iota
	Class
	Tag
)

type Selector struct {
	selector     string
	selectorName string
	selectorType BasicSelector
}

func matchNode(node *Node, basicSelectorName string, basicSelectorType BasicSelector) bool {
	if basicSelectorName == "" {
		return true
	} else if node == nil {
		return false
	}

	switch basicSelectorType {
	case Id:
		idName, _ := node.GetAttribute("id")
		return idName == basicSelectorName
	case Class:
		classList := NewClassList()
		classList.DecodeFrom(node)
		return classList.Contains(basicSelectorName)
	case Tag:
		return node.GetTagName() == basicSelectorName
	}
	return false
}

func NewSelector(selector string) Selector {
	selector = strings.TrimSpace(html.EscapeString(selector))
	selectorStruct := Selector{}
	if len(selector) == 0 || (selector[0] == '.' || selector[0] == '#') && len(selector) <= 1 {
		return selectorStruct
	}

	switch selector[0] {
	case '.':
		selectorStruct.selectorType = Class
	case '#':
		selectorStruct.selectorType = Id
	default:
		selectorStruct.selectorType = Tag
	}

	selectorStruct.selector = strings.ToLower(selector)
	if selectorStruct.selectorType != Tag {
		selectorStruct.selectorName = selector[1:]
	} else {
		selectorStruct.selectorName = selector
	}
	return selectorStruct
}

type Combinator int

const (
	Descendant Combinator = iota
	Child
	NextSibling
	SubsequentSibling
	//if no combinator
	NoneCombinator
)

type CombinatorEl struct {
	Type      Combinator
	Selector1 Selector
	Selector2 Selector
}

func TokenizeSelectorsAndCombinators(selector string) []CombinatorEl {
	iter := func(yield func(string) bool) {
		currentStr := ""
		for _, char := range selector {
			switch char {
			case ' ', '>', '+', '~':
				if !yield(currentStr) || !yield(string(char)){
					return
				}
				currentStr = ""
			default:
				currentStr+=string(char)
			}
		}
		yield(currentStr)
	}

	list := make([]CombinatorEl, 0, 1)
	currentCombinator := *new(CombinatorEl)
	currentCombinator.Selector1 = NewSelector("")
	currentCombinator.Type = NoneCombinator
	for str := range iter {
		if strings.TrimSpace(str) == "" {
			continue
		}

		switch str {
		case "+":
			currentCombinator.Type = NextSibling
		case ">":
			currentCombinator.Type = Child
		case "~":
			currentCombinator.Type = SubsequentSibling
		default:
			newSelector := NewSelector(str)
			currentCombinator.Selector2 = newSelector
			list = append(list, currentCombinator)
			currentCombinator = *new(CombinatorEl)
			currentCombinator.Selector1 = newSelector
		}

	}

	if len(list) == 1 {
		list[0].Type = NoneCombinator
	}

	return list
}

func (ce *CombinatorEl) getMatchingNode(node *Node) *Node {
	switch ce.Type {
	case Descendant:
		return ce.getDescended(node)
	case Child:
		return ce.getDirectChild(node)
	case NextSibling:
		return ce.getNextSibling(node)
	case SubsequentSibling:
		return ce.getSubsequentSibling(node)
	case NoneCombinator:
		if matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) {
			return node
		}
	}
	return nil
}

// isDescended returns wether the given node is a ce.Selector2 and descended of ce.Selector1.
func (ce *CombinatorEl) getDescended(node *Node) *Node {
	if !matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) {
		return nil
	}

	parentNode := node.GetParent()
	for parentNode != nil {
		if matchNode(parentNode, ce.Selector1.selectorName, ce.Selector1.selectorType) {
			return parentNode
		}
		parentNode = parentNode.GetParent()
	}
	return nil
}

// isDirectChild returns whether the given node is a direct child of ce.Selector1 and node is of ce.Selector2
func (ce *CombinatorEl) getDirectChild(node *Node) *Node {
	if node == nil {
		return nil
	}

	if matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) &&
		matchNode(node.GetParent(), ce.Selector1.selectorName, ce.Selector1.selectorType) {
		return node.GetParent()
	}
	return nil
}

// isNextSibling return whether the given node is of ce.Selector2 and next sibling of ce.Selector1
func (ce *CombinatorEl) getNextSibling(node *Node) *Node {
	if node == nil {
		return nil
	}

	if matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) &&
		matchNode(node.GetPreviousNode(), ce.Selector1.selectorName, ce.Selector1.selectorType) {
		return node.GetPreviousNode()
	}
	return nil
}

func (ce *CombinatorEl) getSubsequentSibling(node *Node) *Node {
	if node == nil || !matchNode(node, ce.Selector2.selector, ce.Selector2.selectorType) {
		return nil
	}

	traverser := NewTraverser(node)
	for traverser.GetCurrentNode() != nil {
		if matchNode(traverser.GetCurrentNode(), ce.Selector1.selector, ce.Selector1.selectorType) {
			return traverser.GetCurrentNode()
		}
		traverser.Previous()
	}
	return nil
}
