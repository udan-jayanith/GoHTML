package GoHtml

import (
	"strings"
)

func CreateNode(tagName string) *Node {
	return &Node{
		TagName: strings.ToLower(strings.TrimSpace(tagName)),
		Closed:  true,
	}
}

func CloneNode(node *Node) *Node{
	newNode := *node
	newNode.parentNode = nil
	newNode.PreviousNode = nil
	newNode.NextNode = nil

	return &newNode
}
