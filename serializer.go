package GoHtml

import (
	"fmt"
	"io"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

func wrapAttributeValue(value string) string {
	if isDigit(value) {
		return value
	}

	return `"` + strings.ReplaceAll(value, `"`, "&quot;") + `"`
}

func encodeListAttributes(node *Node) string {
	w := strings.Builder{}
	attributes := node.IterAttributes()
	for a, v := range attributes{
		if strings.TrimSpace(v) == "" {
			w.Write(fmt.Appendf(nil, " %s", a))
		} else {
			w.Write(fmt.Appendf(nil, " %s=%s", a, wrapAttributeValue(v)))
		}
	}
	return w.String()
}

// Encode writes to w encoding of rootNode
func Encode(w io.Writer, rootNode *Node) {
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
			fmt.Fprintf(w, "<%s%s>", tagName, encodeListAttributes(current))
			if current.GetNextNode() != nil {
				stack.Push(stackFrame{node: current.GetNextNode(), openedTag: false})
			}
		} else if !top.openedTag {
			fmt.Fprintf(w, "<%s%s>", tagName, encodeListAttributes(current))
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

// NodeTreeToHTML returns encoding of node-tree as a string.
func NodeTreeToHTML(rootNode *Node) string {
	builder := &strings.Builder{}
	Encode(builder, rootNode)
	return builder.String()
}