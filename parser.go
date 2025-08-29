package GoHtml

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Decode reads from rd and create a node-tree. Then returns the root node and nil.
func Decode(r io.Reader) (*Node, error) {
	t := NewTokenizer(r)
	nodeTreeBuilder := NewNodeTreeBuilder()
	for {
		tt := t.Advanced()
		if tt == html.ErrorToken {
			break
		}

		nodeTreeBuilder.WriteNodeTree(t.GetCurrentNode(), tt)
	}
	return nodeTreeBuilder.GetRootNode(), nil
}

// HTMLToNodeTree return html code as a node-tree. If error were to occur it would be SyntaxError.
func HTMLToNodeTree(html string) (*Node, error) {
	rd := strings.NewReader(html)
	node, err := Decode(rd)
	return node, err
}

