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



func serializeHTMLTag(tag string) (*Node, error) {
	tag = strings.TrimRight(strings.TrimLeft(tag, "<"), ">")
	reg := regexp.MustCompile(`(\w+(?:-\w+)*)\s*(?:=\s*(?:(["'`+"`"+`])(.*?)2|(\S+)))?`).FindAllString(tag, len(tag))
	if reg == nil {
		return CreateNode(""), fmt.Errorf("Invalid html tag")
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

func isReadingTag(strBuf string) bool {
	return regexp.MustCompile(`^<.*`).MatchString(strBuf)
}

func isClosingTag(tag string) bool {
	reg := regexp.MustCompile(`^<\/.*>\s*$`)
	return reg.MatchString(tag)
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
