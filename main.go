/*
A HTML parse and a serializer for Go. GoHTML tries to keep semantic similar to JS-DOM API while trying to keep the API simple by not forcing JS-DOM model into GoHTML. Because of this GoHTML has node tree model. GoHTML tokenizer uses std net/html module for tokenizing in underlining layer. There for it's users responsibility to make sure inputs to GoHTML is UTF-8 encoded. GoHTML allows direct access to the node tree.
*/
package GoHtml

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	SyntaxError error = fmt.Errorf("Syntax error")
)

// CreateNode returns a initialized new node.
func CreateNode(tagName string) *Node {
	return &Node{
		tagName:    strings.ToLower(strings.TrimSpace(tagName)),
		attributes: make(map[string]string),
	}
}

// CreateTextNode returns a new node that represents the given text.
// HTML tags in text get escaped.
func CreateTextNode(text string) *Node {
	textNode := CreateNode("")
	textNode.SetText(text)
	return textNode
}

// DeepCloneNode clones the node without having references to it's original parent node, previous node and next node.
// If node is nil DeepCloneNode returns nil.
func DeepCloneNode(node *Node) *Node {
	if node == nil {
		return node
	}
	attributes := node.attributes

	newNode := Node{
		childNode:  node.GetChildNode(),
		tagName:    node.GetTagName(),
		attributes: attributes,
		text:       node.GetText(),
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