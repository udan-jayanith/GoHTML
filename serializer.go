package GoHtml

import (
	"fmt"
	"io"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	"golang.org/x/net/html"
)

func encodeListAttributes(node *Node) string {
	w := strings.Builder{}
	node.IterateAttributes(func(attribute, value string) {
		if strings.TrimSpace(value) == "" {
			w.Write(fmt.Appendf(nil, " %s", attribute))
		} else {
			w.Write(fmt.Appendf(nil, " %s=%s", attribute, `"`+html.EscapeString(value)+`"`))
		}

	})
	return w.String()
}

// Encode writes to w encoding of the node tree from rootNode.
func Encode(w io.Writer, rootNode *Node) {
	if rootNode == nil {
		return
	}

	type stackFrame struct {
		node         *Node
		isClosingTag bool
	}

	stack := linkedliststack.New()
	stack.Push(stackFrame{
		node: rootNode,
	})

	for stack.Size() > 0 {
		v, _ := stack.Pop()
		currentStackFrame := v.(stackFrame)

		if currentStackFrame.isClosingTag {
			fmt.Fprintf(w, "</%s>", currentStackFrame.node.GetTagName())
			continue
		} else if currentStackFrame.node.IsTextNode() {
			fmt.Fprint(w, html.EscapeString(currentStackFrame.node.GetText()))
		} else {
			fmt.Fprintf(w, "<%s%s>", func() string {
				tagName := currentStackFrame.node.GetTagName()
				tagNameUpperCased := strings.ToUpper(tagName)
				if tagNameUpperCased == DOCTYPEDTD {
					tagName = tagNameUpperCased
				}
				return tagName
			}(), encodeListAttributes(currentStackFrame.node))
		}

		if currentStackFrame.node.GetNextNode() != nil {
			stack.Push(stackFrame{
				node: currentStackFrame.node.GetNextNode(),
			})
		}
		if !IsVoidTag(currentStackFrame.node.GetTagName()) && !currentStackFrame.node.IsTextNode(){
			stack.Push(stackFrame{
				node:         currentStackFrame.node,
				isClosingTag: true,
			})
		}
		if currentStackFrame.node.GetChildNode() != nil {
			stack.Push(stackFrame{
				node: currentStackFrame.node.GetChildNode(),
			})
		}
	}
}

// NodeTreeToHTML returns encoding of node-tree as a string.
func NodeTreeToHTML(rootNode *Node) string {
	builder := &strings.Builder{}
	Encode(builder, rootNode)
	return builder.String()
}
