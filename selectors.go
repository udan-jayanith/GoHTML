package GoHtml

import (
	"strings"
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
	if basicSelectorName == ""{
		return true
	}else if node == nil {
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
	selector = strings.TrimSpace(selector)
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

	selectorStruct.selector = selector
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
	list := make([]CombinatorEl, 0, 1)
	slice := strings.SplitSeq(selector, " ")
	currentCombinator := *new(CombinatorEl)
	currentCombinator.Selector1 = NewSelector("")
	for str := range slice {
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

func (ce *CombinatorEl) IsMatchingNode(node *Node) bool {
	switch ce.Type {
	case Descendant:
		return ce.isDescended(node)
	case Child:
		return ce.isDirectChild(node)
	case NextSibling:
		return ce.isNextSibling(node)
	case SubsequentSibling:
		return ce.isSubsequentSibling(node)
	case NoneCombinator:
		return matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType)
	}
	return false
}

// isDescended returns wether the given node is a ce.Selector2 and descended of ce.Selector1.
func (ce *CombinatorEl) isDescended(node *Node) bool {
	if !matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) {
		return false
	}

	parentNode := node.GetParent()
	for parentNode != nil && !matchNode(parentNode, ce.Selector1.selectorName, ce.Selector1.selectorType) {
		parentNode = parentNode.GetParent()
	}
	return parentNode != nil
}

// isDirectChild returns whether the given node is a direct child of ce.Selector1 and node is of ce.Selector2
func (ce *CombinatorEl) isDirectChild(node *Node) bool {
	if node == nil {
		return false
	}

	return matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) && matchNode(node.GetParent(), ce.Selector1.selectorName, ce.Selector1.selectorType)
}

// isNextSibling return whether the given node is of ce.Selector2 and next sibling of ce.Selector1
func (ce *CombinatorEl) isNextSibling(node *Node) bool {
	if node == nil {
		return false
	}

	return matchNode(node, ce.Selector2.selectorName, ce.Selector2.selectorType) && matchNode(node.GetPreviousNode(), ce.Selector1.selectorName, ce.Selector1.selectorType)
}

func (ce *CombinatorEl) isSubsequentSibling(node *Node) bool {
	if !matchNode(node, ce.Selector2.selector, ce.Selector2.selectorType) {
		return false
	}

	traverser := NewTraverser(node)
	for traverser.GetCurrentNode() != nil && !matchNode(traverser.GetCurrentNode(), ce.Selector1.selector, ce.Selector1.selectorType) {
		traverser.Previous()
	}
	return matchNode(traverser.GetCurrentNode(), ce.Selector1.selector, ce.Selector1.selectorType)
}
