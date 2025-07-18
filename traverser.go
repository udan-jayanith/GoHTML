package GoHtml

type Traverser struct{
	CurrentNode *Node
}

//GetTraverser returns a new traverser that can be used to navigate the node tree.
func GetTraverser(startingNode *Node) Traverser{
	return Traverser{
		CurrentNode: startingNode,
	}
}

//Next returns the node next to current node and change CurrentNode to the new node.
func (t *Traverser) Next() *Node{
	t.CurrentNode = t.CurrentNode.NextNode
	return t.CurrentNode
} 

//Previous returns the previous node and change CurrentNode to the new node.
func (t *Traverser) Previous() *Node{	
	t.CurrentNode = t.CurrentNode.PreviousNode
	return t.CurrentNode
}