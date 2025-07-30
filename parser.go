package GoHtml

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

// Decode reads from rd and create a node-tree. Then returns the root node and an error. If error were to occur it would be SyntaxError.
func Decode(rd io.Reader) (*Node, error) {
	newRd := bufio.NewReader(rd)
	rootNode := CreateNode("")
	currentNode := rootNode
	stack := linkedliststack.New()

	str := ""
	readingQuote := ""
	readingComment := ""
	for {
		byt, err := newRd.ReadByte()
		if err != nil {
			node := rootNode.GetNextNode()
			rootNode.RemoveNode()
			return node, nil
		}
		str += string(byt)

		if readingComment == "" && isStartingComment(currentNode, str) && readingQuote == "" && currentNode.GetTagName() != "script"{
			readingComment = getStartingComment(currentNode, str)
		} else if readingComment != "" && isEndingComment(currentNode, readingComment, str) {
			readingComment = ""
			str = ""
		}

		if readingComment != "" {
			continue
		}

		if isQuote(string(byt)) && (currentNode.GetTagName() == "script" || currentNode.GetTagName() == "style" || isReadingTag(str)) &&
			readingQuote == "" {
			readingQuote = string(byt)
			continue
		} else if readingQuote == string(byt) {
			readingQuote = ""
		}

		if readingQuote != "" {
			continue
		}

		if regexp.MustCompile(`^\s*<\/.*>\s*$`).MatchString(str) {
			//closing tag
			str = ""
			currentNode, err = getFirstOpenNode(currentNode, stack)
			if err != nil {
				return currentNode, err
			}
		} else if regexp.MustCompile(`^\s*<.*>\s*$`).MatchString(str) {
			//opening and void tags
			if isJsComment(str) && currentNode.GetTagName() == "script"{
				str = ""
				continue
			}

			node, err := serializeHTMLTag(str)
			if err != nil {
				node := rootNode.GetNextNode()
				rootNode.RemoveNode()
				return node, err
			}
			str = ""

			if isOpen(currentNode, stack) {
				currentNode.AppendChild(node)
			} else {
				currentNode.Append(node)
			}
			currentNode = node
			if !isSelfClosingNode(node) {
				stack.Push(currentNode)
			}
		} else if string(byt) == "<" && !regexp.MustCompile(`^\s*<$`).MatchString(str) {
			// text
			str = str[:len(str)-1]
			node := serializeTextNode(str)
			str = "<"

			if isOpen(currentNode, stack) {
				currentNode.AppendChild(node)
			} else {
				currentNode.Append(node)
			}
			currentNode = node
		}
	}
}

func isStartingComment(currentNode *Node, last4Char string) bool {
	if currentNode.GetTagName() == "script" {
		return regexp.MustCompile(`//$`).MatchString(last4Char) || regexp.MustCompile(`/\*$`).MatchString(last4Char)
	} else if currentNode.GetTagName() == "style" {
		return regexp.MustCompile(`/\*$`).MatchString(last4Char)
	}

	return regexp.MustCompile(`<!--$`).MatchString(last4Char)
}

func getStartingComment(currentNode *Node, last4Char string) string {
	if currentNode.GetTagName() == "script" {
		if regexp.MustCompile(`//$`).MatchString(last4Char) {
			return "//"
		}
		return regexp.MustCompile(`/\*$`).FindString(last4Char)
	} else if currentNode.GetTagName() == "style" {
		return regexp.MustCompile(`/\*$`).FindString(last4Char)
	}

	return regexp.MustCompile(`<!--$`).FindString(last4Char)
}

func isEndingComment(currentNode *Node, startingComment string, last4Char string) bool {
	if currentNode.GetTagName() == "script" {
		return (regexp.MustCompile(`\n$`).MatchString(last4Char) && startingComment == "//") || (regexp.MustCompile(`\*/$`).MatchString(last4Char) && startingComment == "/*")
	} else if currentNode.GetTagName() == "style" {
		return regexp.MustCompile(`\*/$`).MatchString(last4Char) && startingComment == "/*"
	}

	return regexp.MustCompile(`-->$`).MatchString(last4Char) && startingComment == "<!--"
}

func isJsComment(str string) bool{
	jsCommentReg := regexp.MustCompile(`^\s*<!--.*-->\s*$`)
	return jsCommentReg.MatchString(str)
}

func getFirstOpenNode(currentNode *Node, stack *linkedliststack.Stack) (*Node, error) {
	traverser := GetTraverser(currentNode)
	for traverser.GetCurrentNode() != nil {
		n, ok := stack.Peek()
		if !ok {
			return traverser.GetCurrentNode(), SyntaxError
		}
		node := n.(*Node)

		if traverser.GetCurrentNode() == node {
			stack.Pop()
			return node, nil
		}

		if traverser.GetCurrentNode().GetPreviousNode() == nil {
			traverser.SetCurrentNodeTo(traverser.GetCurrentNode().GetParent())
		} else {
			traverser.Previous()
		}
	}

	return traverser.GetCurrentNode(), SyntaxError
}

func isOpen(currentNode *Node, stack *linkedliststack.Stack) bool {
	if stack.Size() < 1 || isSelfClosingNode(currentNode) {
		return false
	}

	n, _ := stack.Peek()
	node := n.(*Node)
	return node == currentNode
}

func isSelfClosingNode(node *Node) bool {
	return node.GetTagName() == "" || IsVoidTag(node.GetTagName())
}

func serializeHTMLTag(tag string) (*Node, error) {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimRight(strings.TrimLeft(tag, "<"), ">")
	node := CreateNode("")

	//extract the html tag name
	tagNameRegex := regexp.MustCompile(`^\s*([\w\-_!]*)`)
	tagName := tagNameRegex.FindString(tag)
	if tagName == "" {
		return node, SyntaxError
	}
	node.SetTagName(strings.TrimSpace(tagName))

	//Cut the tag name from tag.
	afterTagName := regexp.MustCompile(`^\s*[\w\-_]*\s*(.*)`)
	subMatch := afterTagName.FindStringSubmatch(tag)
	if len(subMatch) > 1 {
		if strings.TrimSpace(getLeftMostString(subMatch)) == "" {
			return node, nil
		}
		tag = getLeftMostString(subMatch)
	} else {
		return node, SyntaxError
	}
	if tag == "" || strings.TrimSpace(getLeftMostString(afterTagName.FindStringSubmatch(tag))) == strings.TrimSpace(tagName){
		return node, nil
	}

	for {
		if strings.TrimSpace(tag) == "" {
			return node, nil
		}

		//This parses attribute name.
		attributeNameReg := regexp.MustCompile(`^\s*([\w\-_]*)\s*`)
		subMatch = attributeNameReg.FindStringSubmatch(tag)
		attributeName := ""
		if len(subMatch) > 0 {
			attributeName = getLeftMostString(subMatch)
		} else {
			return node, SyntaxError
		}

		//This remove attribute name from the tag.
		afterAttributeName := regexp.MustCompile(`^\s*[\w\-_!]*\s*(.*)`)
		subMatch = afterAttributeName.FindStringSubmatch(tag)
		if len(subMatch) >= 1 {
			if strings.TrimSpace(getLeftMostString(subMatch)) == "" {
				return node, nil
			}
			tag = getLeftMostString(subMatch)
		} else {
			return node, SyntaxError
		}

		isDefinedValueReg := regexp.MustCompile(`^\s*(=).*`)
		if !isDefinedValueReg.MatchString(tag) {
			node.SetAttribute(attributeName, "")
			if strings.TrimSpace(attributeName) == strings.TrimSpace(getLeftMostString(attributeNameReg.FindStringSubmatch(tag))){
				return node, nil
			}
			continue
		}

		afterEqualSignReg := regexp.MustCompile(`^\s*=(.*)`)
		subMatch = afterEqualSignReg.FindStringSubmatch(tag)
		if len(subMatch) >= 1 {
			tag = getLeftMostString(subMatch)
		} else {
			return node, SyntaxError
		}

		definedValueReg := regexp.MustCompile(`\s*(\s*('.*?'|(".*?")|[\d+\.?\d+]*).*).*`)
		subMatch = definedValueReg.FindStringSubmatch(tag)
		if len(subMatch) >= 1 {
			node.SetAttribute(attributeName, escapeQuotes(getLeftMostString(subMatch)))
		} else {
			return node, SyntaxError
		}

		afterDefinedValueReg := regexp.MustCompile(`\s*('.*?'|(".*?")|[\d+\.?\d+]*)\s*(.*)`)
		subMatch = afterDefinedValueReg.FindStringSubmatch(tag)
		if len(subMatch) >= 1 {
			if strings.TrimSpace(getLeftMostString(subMatch)) == "" || strings.TrimSpace(getLeftMostString(subMatch))  == strings.TrimSpace(tag){
				return node, nil
			}
			tag = getLeftMostString(subMatch)
		} else {
			return node, SyntaxError
		}

	}
}

func serializeTextNode(s string) *Node {
	node := CreateTextNode(s)
	return node
}

func isReadingTag(strBuf string) bool {
	return regexp.MustCompile(`^<.*`).MatchString(strBuf)
}

// HTMLToNodeTree return html code as a node-tree. If error were to occur it would be SyntaxError.
func HTMLToNodeTree(html string) (*Node, error) {
	rd := strings.NewReader(html)
	node, err := Decode(rd)
	return node, err
}

func escapeQuotes(str string) string {
	escapeQuotesReg := regexp.MustCompile(`^\s*('(.*)'|"(.*)"|([\d+\.]*)|.*)\s*$`)
	matches := escapeQuotesReg.FindStringSubmatch(str)
	for i := len(matches) - 1; i >= 0; i-- {
		if strings.TrimSpace(matches[i]) != "" {
			return matches[i]
		}
	}
	return ""
}

func getLeftMostString(slice []string) string {
	for i := len(slice) - 1; i >= 0; i-- {
		if strings.TrimSpace(slice[i]) != "" {
			return slice[i]
		}
	}
	return ""
}
