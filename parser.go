package GoHtml

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

var(
	closingTagReg = regexp.MustCompile(`^\s*<\/.*>\s*$`)
	openingTagReg =  regexp.MustCompile(`^\s*<.*>\s*$`)
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

		if readingComment == "" && isStartingComment(currentNode, str) && readingQuote == "" {
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

		if closingTagReg.MatchString(str) {
			//closing tag
			str = ""
			currentNode, err = getFirstOpenNode(currentNode, stack)
			if err != nil {
				return currentNode, err
			}
		} else if openingTagReg.MatchString(str) {
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
		} else if string(byt) == "<" && strings.TrimSpace(str) != "<" {
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

var (
	startingBlockCommentReg = regexp.MustCompile(`\s*\/\*\s*$`)
	endingBlockCommentReg = regexp.MustCompile(`\s*\*\/\s*$`)
	doubleSlashReg = regexp.MustCompile(`//$`)
	htmlCommentStarterReg = regexp.MustCompile(`<!--$`)
)

func isStartingComment(currentNode *Node, str string) bool {
	if currentNode.GetTagName() == "script" {
		return doubleSlashReg.MatchString(str) || startingBlockCommentReg.MatchString(str)
	} else if currentNode.GetTagName() == "style" {
		return startingBlockCommentReg.MatchString(str)
	}

	return htmlCommentStarterReg.MatchString(str)
}

var(
	endingNewLineReg = regexp.MustCompile(`\n$`)
	htmlCommentEndReg = regexp.MustCompile(`-->$`)
)

func isEndingComment(currentNode *Node, startingComment string, str string) bool {
	if currentNode.GetTagName() == "script" {
		return (endingNewLineReg.MatchString(str) && startingComment == "//") || (endingBlockCommentReg.MatchString(str) && startingComment == "/*")
	} else if currentNode.GetTagName() == "style" {
		return endingBlockCommentReg.MatchString(str) && startingComment == "/*"
	}

	return htmlCommentEndReg.MatchString(str) && startingComment == "<!--"
}

func getStartingComment(currentNode *Node, str string) string {
	if currentNode.GetTagName() == "script" {
		if doubleSlashReg.MatchString(str) {
			return "//"
		}
		return getRightMostString(startingBlockCommentReg.FindStringSubmatch(str))
	} else if currentNode.GetTagName() == "style" {
		return getRightMostString(startingBlockCommentReg.FindStringSubmatch(str))
	}

	return getRightMostString(htmlCommentStarterReg.FindStringSubmatch(str))
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

var (
	tagNameRegex          = regexp.MustCompile(`^\s*([\w\-_!]*)`)
	afterTagNameReg       = regexp.MustCompile(`^\s*[\w\-_]*\s*(.*)`)
	attributeNameReg      = regexp.MustCompile(`^\s*([\w\-_!]*)\s*`)
	afterAttributeNameReg = regexp.MustCompile(`^\s*[\w\-_!]*\s*(.*)`)
	isDefinedValueReg     = regexp.MustCompile(`^\s*=.*`)
	afterEqualSignReg     = regexp.MustCompile(`^\s*=(.*)`)
	definedValueReg       = regexp.MustCompile(`\s*(\s*('.*?'|".*?"|\s*[\S]*).*).*`)
	afterDefinedValueReg  = regexp.MustCompile(`\s*('.*?'|".*?"|\s*[\S]*)\s*(.*)\s*`)
)

func serializeHTMLTag(tag string) (*Node, error) {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimRight(strings.TrimRight(strings.TrimLeft(tag, "<"), ">"), `/`)
	node := CreateNode("")

	//extract the html tag name
	tagName := tagNameRegex.FindString(tag)
	if tagName == "" {
		return node, SyntaxError
	}
	node.SetTagName(strings.TrimSpace(tagName))

	//Cut the tag name from tag.
	tag = strings.TrimSpace(getRightMostString(afterTagNameReg.FindStringSubmatch(tag)))
	if strings.TrimSpace(tag) == "" || tag == strings.TrimSpace(tagName) {
		return node, nil
	}

	for {
		if tag == "" {
			return node, nil
		}

		//This parses attribute name.
		attributeName := strings.TrimSpace(getRightMostString(attributeNameReg.FindStringSubmatch(tag)))
		if attributeName == "" {
			return node, SyntaxError
		}

		//This removes attribute name from the tag.
		tag = strings.TrimSpace(getRightMostString(afterAttributeNameReg.FindStringSubmatch(tag)))
		if tag == "" {
			return node, nil
		}

		if !isDefinedValueReg.MatchString(tag) {
			node.SetAttribute(attributeName, "")
			if attributeName == strings.TrimSpace(tag){
				return node, nil
			}
			continue
		}

		tag = strings.TrimSpace(getRightMostString(afterEqualSignReg.FindStringSubmatch(tag)))
		if tag == "" {
			return node, SyntaxError
		}

		attributeValue := strings.TrimSpace(getRightMostString(definedValueReg.FindStringSubmatch(tag)))
		node.SetAttribute(attributeName, escapeQuotes(attributeValue))

		tag = strings.TrimSpace(getRightMostString(afterDefinedValueReg.FindStringSubmatch(tag)))
		if tag == "" || attributeValue == tag{
			return node, nil
		}

	}
}

func serializeTextNode(s string) *Node {
	node := CreateTextNode(s)
	return node
}

var (
	firstCharLesserReg *regexp.Regexp = regexp.MustCompile(`^<.*`)
)

func isReadingTag(strBuf string) bool {
	return firstCharLesserReg.MatchString(strBuf)
}

// HTMLToNodeTree return html code as a node-tree. If error were to occur it would be SyntaxError.
func HTMLToNodeTree(html string) (*Node, error) {
	rd := strings.NewReader(html)
	node, err := Decode(rd)
	return node, err
}

var (
	escapeQuotesReg *regexp.Regexp = regexp.MustCompile(`^\s*('(.*)'|"(.*)"|([\d+\.]*)|.*)\s*$`)
)

func escapeQuotes(str string) string {
	matches := escapeQuotesReg.FindStringSubmatch(str)
	for i := len(matches) - 1; i >= 0; i-- {
		if strings.TrimSpace(matches[i]) != "" {
			return matches[i]
		}
	}
	return ""
}

func getRightMostString(slice []string) string {
	for i := len(slice) - 1; i >= 0; i-- {
		if strings.TrimSpace(slice[i]) != "" {
			return slice[i]
		}
	}
	return ""
}
