package GoHtml

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

var (
	SyntaxError error = fmt.Errorf("Syntax error")
)

// Given the stack and the currentNode check wether the current node is closed or open. If current node is closed append the node as a next node otherwise append the node as a child node.
// If the tag is a closing tag close the first unclosed tag and add that tag to the stack. Pop nodes from the stack until first unclosed tag and set current node as that closed node.
// if the s is a closing tag get the first unclosed tag and add it to the stack.
func DecodeToNodeTree(rd io.Reader) (*Node, error) {
	newRd := bufio.NewReader(rd)
	rootNode := CreateNode("")
	currentNode := rootNode
	stack := linkedliststack.New()
	markAsClosed(currentNode, stack)

	str := ""
	readingQuote := ""
	for {
		byt, err := newRd.ReadByte()
		if err != nil {
			return rootNode, nil
		}
		str += string(byt)

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

			markAsClosed(currentNode, stack)
			node, err := getFirstUnclosedNode(currentNode, stack)
			if err != nil {
				return rootNode, err
			}
			currentNode = node
			str = ""
		} else if regexp.MustCompile(`^\s*<.*>\s*$`).MatchString(str) {
			//opening tag
			node, err := serializeHTMLTag(str)
			if err != nil {
				return rootNode, err
			}
			str = ""

			if isClosed(currentNode, stack) {
				currentNode.Append(node)
			} else {
				currentNode.AppendChild(node)
			}
			currentNode = node
		} else if string(byt) == "<" && !regexp.MustCompile(`^\s*<$`).MatchString(str) {
			// text
			str = str[:len(str)-1]
			node := serializeTextNode(str)
			str = "<"

			if isClosed(currentNode, stack) {
				currentNode.Append(node)
			} else {
				currentNode.AppendChild(node)
			}
			currentNode = node
		}
	}
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
			node.SetAttribute(strings.TrimSpace(v), "true")
		}
	}
	return node, nil
}

func serializeTextNode(s string) *Node {
	node := CreateTextNode(s)
	return node
}

func getFirstUnclosedNode(currentNode *Node, stack *linkedliststack.Stack) (*Node, error) {
	traverser := GetTraverser(currentNode)
	for traverser.GetCurrentNode() != nil {
		if !isClosed(traverser.GetCurrentNode(), stack) {
			return traverser.GetCurrentNode(), nil
		}
		
		// if node is a text node or a void node doesn't have to pop from the stack
		if traverser.GetCurrentNode().GetTagName() != "" && !IsVoidTag(traverser.GetCurrentNode().GetTagName()){
			_, ok := stack.Pop()
			if !ok {
				return traverser.GetCurrentNode(), SyntaxError
			}
		}

		if traverser.GetCurrentNode().GetPreviousNode() == nil {
			traverser.SetCurrentNodeTo(traverser.GetCurrentNode().GetParent())
		} else {
			traverser.Previous()
		}
	}

	return traverser.GetCurrentNode(), SyntaxError
}

func isReadingTag(strBuf string) bool {
	return regexp.MustCompile(`^<.*`).MatchString(strBuf)
}

func isClosed(currentNode *Node, stack *linkedliststack.Stack) bool {
	if currentNode.GetTagName() == "" || IsVoidTag(currentNode.GetTagName()) {
		return true
	}else if stack.Size() < 1 {
		return false
	}

	m, _ := stack.Peek()
	node := m.(*Node)
	return currentNode == node
}

func markAsClosed(currentNode *Node, stack *linkedliststack.Stack) {
	if currentNode.GetTagName() == "" || IsVoidTag(currentNode.GetTagName()) {
		return
	}
	stack.Push(currentNode)
}

func isQuote(chr string) bool {
	return chr == `"` || chr == `'` || chr == "`"
}

func isDigit(value string) bool {
	reg := regexp.MustCompile(`^[\d\.]+$`)
	return reg.Match([]byte(value))
}

func wrapAttributeValue(value string) string {
	if isDigit(value) {
		return value
	}

	return `"` + strings.ReplaceAll(value, `"`, "&quot;") + `"`
}

func encodeListAttributes(node *Node) string {
	w := strings.Builder{}
	node.IterateAttributes(func(attribute, value string) {
		w.Write(fmt.Appendf(nil, " %s=%s", attribute, wrapAttributeValue(value)))
	})
	return w.String()
}

func EncodeToHTML(w io.Writer, rootNode *Node) {
	type stackFrame struct {
		node      *Node
		openedTag bool
	}

	stack := linkedliststack.New()
	stack.Push(stackFrame{node: rootNode, openedTag: false})

	for stack.Size() > 0 {
		t, _ := stack.Pop()
		top := t.(stackFrame)
		current := top.node

		if current == nil {
			continue
		}

		tagName := current.GetTagName()
		if tagName == "" {
			w.Write([]byte(current.GetText()))
		} else if IsVoidTag(tagName) {
			fmt.Fprintf(w, "<%s %s>", tagName, encodeListAttributes(current))
			if current.GetNextNode() != nil {
				stack.Push(stackFrame{node: current.GetNextNode(), openedTag: false})
			}
		} else if !top.openedTag {
			fmt.Fprintf(w, "<%s %s>", tagName, encodeListAttributes(current))
			stack.Push(stackFrame{node: current, openedTag: true})

			if current.GetChildNode() != nil {
				stack.Push(stackFrame{node: current.GetChildNode(), openedTag: false})
			}
		} else {
			fmt.Fprintf(w, "</%s>", tagName)
			if current.GetNextNode() != nil {
				stack.Push(stackFrame{node: current.GetNextNode(), openedTag: false})
			}
		}
	}
}
