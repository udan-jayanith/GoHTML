package GoHtml

import (
	"bufio"
	"fmt"
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

		last4Char := getLast4Char(str)
		if readingComment == "" && isStartingComment(currentNode, last4Char) {
			readingComment = getStartingComment(currentNode, last4Char)
		} else if readingComment != "" && isEndingComment(currentNode, readingComment, last4Char) {
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
			fmt.Println(str)
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

func getLast4Char(str string) string {
	return regexp.MustCompile(`.{0,4}$`).FindString(str)
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
	tag = strings.TrimRight(strings.TrimLeft(tag, "<"), ">")
	reg := regexp.MustCompile(`([\w!]+(?:-\w+)*)\s*(?:=\s*(?:(["'`+"`"+`])(.*?)2|(\S+)))?`).FindAllString(tag, len(tag))
	if reg == nil {
		return CreateNode(""), SyntaxError
	}

	node := CreateNode(reg[0])
	if len(reg) <= 1 {
		return node, nil
	}
	reg = reg[1:]
	for _, v := range reg {
		if regexp.MustCompile(`^\s*.+\s*=\s*\d+\s*$`).MatchString(v) {
			s := regexp.MustCompile(`^\s*(.+)\s*=\s*(\d+)\s*$`).FindAllStringSubmatch(v, 2)
			node.SetAttribute(s[0][1], s[0][2])
		} else if regexp.MustCompile(`\w+\s*(=)\s*.+`).Match([]byte(v)) {
			s := regexp.MustCompile(`^\s*(.+)\s*=\s*['"](.+)['"]\s*$`).FindAllStringSubmatch(v, 2)
			if len(s) < 1 {
				continue
			}
			node.SetAttribute(s[0][1], s[0][2])
		} else {
			node.SetAttribute(strings.TrimSpace(v), "")
		}
	}
	return node, nil
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
