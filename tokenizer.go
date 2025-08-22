package GoHtml

import (
	"io"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	"golang.org/x/net/html"
)

// Tokenizer contains a *html.Tokenizer.
type Tokenizer struct {
	z *html.Tokenizer
}

// NewTokenizer returns a new Tokenizer.
func NewTokenizer(r io.Reader) Tokenizer {
	return Tokenizer{
		z: html.NewTokenizer(r),
	}
}

// Advanced scans the next token and returns its type.
func (t *Tokenizer) Advanced() html.TokenType {
	return t.z.Next()
}

// CurrentNode returns the current node.
// Returned value can be nil regardless of tt.
func (t *Tokenizer) GetCurrentNode() *Node {
	currentToken := t.z.Token()
	if strings.TrimSpace(currentToken.Data) == "" {
		return nil
	}

	// token data depend on the token type.
	switch currentToken.Type {
	case html.DoctypeToken, html.StartTagToken, html.SelfClosingTagToken, html.TextToken:
		var node *Node
		switch currentToken.Type {
		case html.TextToken:
			node = CreateTextNode(currentToken.Data)
		case html.DoctypeToken:
			node = CreateNode(DOCTYPEDTD)
			node.SetAttribute(currentToken.Data, "")
		default:
			node = CreateNode(currentToken.Data)
			for _, v := range currentToken.Attr {
				node.SetAttribute(v.Key, v.Val)
			}
		}
		return node
	}
	return nil
}

// NodeTreeBuilder is used to build a node tree given a node and it's type.
type NodeTreeBuilder struct {
	rootNode    *Node
	stack       *linkedliststack.Stack
	currentNode *Node
}

// NewNodeTreeBuilder returns a new NodeTreeBuilder.
func NewNodeTreeBuilder() NodeTreeBuilder {
	rootNode := CreateTextNode("")
	return NodeTreeBuilder{
		rootNode:    rootNode,
		currentNode: rootNode,
		stack:       linkedliststack.New(),
	}
}

// WriteNodeTree append the node given html.TokenType
func (ntb *NodeTreeBuilder) WriteNodeTree(node *Node, tt html.TokenType) {
	switch tt {
	case html.EndTagToken:
		val, ok := ntb.stack.Pop()
		if !ok || val == nil {
			return
		}
		ntb.currentNode = val.(*Node)
	case html.DoctypeToken, html.StartTagToken, html.SelfClosingTagToken, html.TextToken:
		if node == nil {
			return
		}

		if isTopNode(ntb.currentNode, ntb.stack) {
			ntb.currentNode.AppendChild(node)
		} else {
			ntb.currentNode.Append(node)
		}

		if !node.IsTextNode() && !IsVoidTag(node.GetTagName()) {
			ntb.stack.Push(node)
		}
		ntb.currentNode = node
	}
}

// GetRootNode returns the root node of the accumulated node tree and resets the NodeTreeBuilder.
func (ntb *NodeTreeBuilder) GetRootNode() *Node {
	node := ntb.rootNode.GetNextNode()
	ntb.rootNode.RemoveNode()

	rootNode := CreateTextNode("")
	ntb.rootNode = rootNode
	ntb.currentNode = rootNode
	ntb.stack = linkedliststack.New()

	return node
}

func isTopNode(node *Node, stack *linkedliststack.Stack) bool {
	val, ok := stack.Peek()
	if !ok || val == nil {
		return false
	}

	topNode := val.(*Node)
	return topNode == node
}