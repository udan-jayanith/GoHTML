package GoHtml

import (
	"io"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	"golang.org/x/net/html"
)

// Decode reads from rd and create a node-tree. Then returns the root node and an error. If error were to occur it would be SyntaxError.
func Decode(r io.Reader) (*Node, error) {
	rootNode := CreateTextNode("")
	stack := linkedliststack.New()
	currentNode := rootNode

	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		currentToken := z.Token()
		if strings.TrimSpace(currentToken.Data) == "" {
			continue
		}

		// token data depend on the token type.
		switch currentToken.Type {
		case html.EndTagToken:
			val, ok := stack.Pop()
			if !ok || val == nil {
				continue
			}
			currentNode = val.(*Node)
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

			if IsTopNode(currentNode, stack){
				currentNode.AppendChild(node)
			}else{
				currentNode.Append(node)
			}

			if !node.IsTextNode() && !IsVoidTag(node.GetTagName()){
				stack.Push(node)
			}
			currentNode = node
		}
	}

	node := rootNode.GetNextNode()
	rootNode.RemoveNode()

	return node, nil
}

// HTMLToNodeTree return html code as a node-tree. If error were to occur it would be SyntaxError.
func HTMLToNodeTree(html string) (*Node, error) {
	rd := strings.NewReader(html)
	node, err := Decode(rd)
	return node, err
}

func IsTopNode(node *Node, stack *linkedliststack.Stack) bool {
	val, ok := stack.Peek()
	if !ok || val == nil {
		return false
	}

	topNode := val.(*Node)
	return topNode == node
}
