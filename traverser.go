package GoHtml

type Traverser struct{
	CurrentNode *Node
}

func GetTraverser(startingNode *Node) Traverser{
	return Traverser{
		CurrentNode: startingNode,
	}
}

//Next returns the node next to current node.
func (t *Traverser) Next() *Node{
	t.CurrentNode = t.CurrentNode.NextNode
	return t.CurrentNode
} 

//Previous return the previous node.
func (t *Traverser) Previous() *Node{	
	t.CurrentNode = t.CurrentNode.PreviousNode
	return t.CurrentNode
}