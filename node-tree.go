package GoHtml


type Node struct {
	NextNode     *Node
	PreviousNode *Node
	ChildNode    *Node
	parentNode   *Node

	TagName    string
	Attributes map[string]string
	Closed     bool
	Text       string
	//RWMutex
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

func (node *Node) GetLastNode() *Node{
	traverser := GetTraverser(node)
	for traverser.CurrentNode.NextNode != nil {
		traverser.Next()
	}
	return traverser.CurrentNode
}

func (node *Node) GetFirstNode() *Node{
	traverser := GetTraverser(node)
	for traverser.CurrentNode.PreviousNode != nil{
		traverser.Previous()
	}
	return traverser.CurrentNode
}

func (node *Node) AppendText(text string){
	textNode := CreateNode("")
	textNode.Text = text

	node.GetLastNode().NextNode = textNode
}

func (node *Node) GetInnerText() string{
	text := ""
	traverser := GetTraverser(node.ChildNode)
	for traverser.CurrentNode != nil{
		text += traverser.CurrentNode.Text
		traverser.Next()
	}

	return text
}
