/*
A powerful and comprehensive HTML parser and DOM manipulation library for Go,
bringing JavaScript-like DOM operations to the Go ecosystem.
*/
package GoHtml

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	SyntaxError error = fmt.Errorf("Syntax error")
)

// CreateNode returns a initialized new node.
func CreateNode(tagName string) *Node {
	return &Node{
		tagName:    strings.ToLower(strings.TrimSpace(tagName)),
		rwMutex:    sync.Mutex{},
		attributes: make(map[string]string),
	}
}

// CreateTextNode returns a new node that represents the given text.
func CreateTextNode(text string) *Node {
	textNode := CreateNode("")
	textNode.SetText(escapeHTML(text))
	return textNode
}

// DeepCloneNode clones the node without having references to it's original parent node, previous node and next node.
// If node is nil DeepCloneNode returns nil.
func DeepCloneNode(node *Node) *Node {
	if node == nil {
		return node
	}

	node.rwMutex.Lock()
	attributes := node.attributes
	node.rwMutex.Unlock()

	newNode := Node{
		childNode:  node.GetChildNode(),
		tagName:    node.GetTagName(),
		attributes: attributes,
		text:       node.GetText(),

		rwMutex: sync.Mutex{},
	}

	return &newNode
}

// CloneNode copy the node. But have one way connections to it's parent, next and previous nodes.
// If node is nil CloneNode returns nil.
func CloneNode(node *Node) *Node {
	if node == nil {
		return nil
	}
	newNode := DeepCloneNode(node)
	newNode.setParentNode(node.getParentNode())
	newNode.SetPreviousNode(node.GetPreviousNode())
	newNode.SetNextNode(node.GetNextNode())

	return newNode
}

func isQuote(chr string) bool {
	return chr == `"` || chr == `'` || chr == "`"
}

func isDigit(value string) bool {
	reg := regexp.MustCompile(`^[\d\.]+$`)
	return reg.Match([]byte(value))
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}
