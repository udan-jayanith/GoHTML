package GoHtml

import (
	"strings"
	"sync"
)

//CreateNode returns a initialized new node.
func CreateNode(tagName string) *Node {
	return &Node{
		tagName: strings.ToLower(strings.TrimSpace(tagName)),
		rwMutex: sync.Mutex{},
		closed:  true,
	}
}

//CloneNode clones the node.
func CloneNode(node *Node) *Node{
	newNode := Node{
		childNode: node.GetChildNode(),
		tagName: node.GetTagName(),
		attributes: node.GetAttributes(),
		text: node.GetText(),
		closed: true,

		rwMutex: sync.Mutex{},
	}

	return &newNode
}
