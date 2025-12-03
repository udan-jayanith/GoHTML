package GoHtml

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

var (
	NoNodesFound error = errors.New("No nodes found in the node tree")
)

// Decode reads from rd and create a node-tree. Then returns the root node and nil.
// If no nodes are found Decode returns a error.
// r must no be nil.
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
	rootNode := nodeTreeBuilder.GetRootNode()
	if rootNode == nil {
		return nil, NoNodesFound
	}
	return rootNode, nil
}

// HTMLToNodeTree return html code as a node-tree. If error were to occur it would be SyntaxError.
func HTMLToNodeTree(html string) (*Node, error) {
	rd := strings.NewReader(html)
	node, err := Decode(rd)
	return node, err
}
