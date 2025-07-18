package GoHtml

type Node struct {
	NextNode     *Node
	PreviousNode *Node
	ChildNodes   *Node
	parentNode   *Node

	TagName    string
	Attributes map[string]string
	Closed     bool
	text       string
}

func (node *Node) AppendChild(childNode *Node) {
	if node.ChildNodes == nil {
		node.ChildNodes = childNode
		return
	}

	traverser := GetTraverser(node.ChildNodes)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	traverser.CurrentNode.NextNode = childNode
}
