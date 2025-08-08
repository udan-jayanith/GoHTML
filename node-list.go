package GoHtml

import (
	"container/list"
	"iter"
)

//NodeList can store nodes by appended order.
type NodeList struct {
	list      *list.List
	currentEl *list.Element
}

// New returns an initialized node list.
func NewNodeList() NodeList {
	return NodeList{
		list:      list.New(),
		currentEl: nil,
	}
}

// Len returns the number of node in the list. The complexity is O(1).
func (nl *NodeList) Len() int {
	return nl.list.Len()
}

// Next advanced to the next node and returns that node.
func (nl *NodeList) Next() *Node {
	if nl.list.Len() == 0 {
		return nil
	}else if nl.currentEl == nil && nl.list.Len() > 0{
		nl.currentEl = nl.list.Front()
	} else {
		if nl.currentEl.Next() == nil{
			return nil
		}
		nl.currentEl = nl.currentEl.Next()
	}
	return nl.currentEl.Value.(*Node)
}

// Previous advanced to the previous node and return that node.
func (nl *NodeList) Previous() *Node {
	if nl.currentEl == nil {
		nl.currentEl = nl.list.Front()
	} else {
		if nl.currentEl.Prev() == nil {
			return nil
		}
		nl.currentEl = nl.currentEl.Prev()
	}
	return nl.currentEl.Value.(*Node)
}

// Back returns the last node of list or nil if the list is empty.
func (nl *NodeList) Back() *Node {
	if nl.list.Back() == nil {
		return nil
	}
	return nl.list.Back().Value.(*Node)
}

// Front returns the first node of list or nil if the list is empty.
func (nl *NodeList) Front() *Node {
	if nl.list.Front() == nil {
		return nil
	}
	return nl.list.Front().Value.(*Node)
}

// Append append a node to the back of the list.
func (nl *NodeList) Append(node *Node) {
	if node == nil{
		return
	}

	nl.list.PushBack(node)
}

// IterNodeList returns a iterator over the node list.
func (nl *NodeList) IterNodeList() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		nodeList := NewNodeList()
		nodeList.list = nl.list

		nextNode := nodeList.Next()
		for nextNode != nil{
			 if !yield(nextNode) {
				return
			}
			nextNode = nodeList.Next()
		}
	}
}
