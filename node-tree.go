package GoHtml

type Node struct {
	NextNode     *Node
	PreviousNode *Node
	ChildNode    *Node
	parentNode   *Node

	TagName    string
	Attributes map[string]string
	Closed     bool
	text       string
}

func (node *Node) AppendChild(childNode *Node) {
	if node.ChildNode == nil {
		node.ChildNode = childNode
		childNode.parentNode = node
		return
	}

	traverser := GetTraverser(node.ChildNode)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	traverser.CurrentNode.NextNode = childNode
	childNode.PreviousNode = traverser.CurrentNode
}

func (node *Node) Append(newNode *Node) {
	traverser := GetTraverser(node)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	newNode.PreviousNode = traverser.CurrentNode
	traverser.CurrentNode.NextNode = newNode
}

func (node *Node) GetParent() *Node {
	traverser := GetTraverser(node)
	for traverser.CurrentNode.parentNode == nil && traverser.CurrentNode.PreviousNode != nil {
		traverser.Previous()
	}
	return traverser.CurrentNode.parentNode
}

//AppendText
//GetInnerText
//Content
//LastNode
//FirstNode