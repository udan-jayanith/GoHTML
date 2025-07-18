package GoHtml

import (
	"strings"
	"sync"
)

//CreateNode returns a initialized new node.
func CreateNode(tagName string) *Node {
	return &Node{
		TagName: strings.ToLower(strings.TrimSpace(tagName)),
		RWMutex: sync.Mutex{},
		AttributesMutex: sync.Mutex{},
		Closed:  true,
	}
}

//CloneNode clones the node.
func CloneNode(node *Node) *Node{
	newNode := Node{
		ChildNode: node.ChildNode,
		TagName: node.TagName,
		Attributes: node.Attributes,
		Text: node.Text,
		Closed: true,

		RWMutex: sync.Mutex{},
		AttributesMutex: sync.Mutex{},
	}

	return &newNode
}
