package GoHtml

import (
	"container/list"
)

type NodeList struct {
	list      *list.List
	currentEl *list.Element
}

// New returns an initialized node list.
func newNodeList() NodeList {
	return *new(NodeList)
}

// Len returns the number of node in the list. The complexity is O(1).
func (nl *NodeList) Len() int {
	return nl.list.Len()
}

// Next returns the next list node or nil.
func (nl *NodeList) Next() *Node {
	if nl.currentEl == nil {
		nl.currentEl = nl.list.Front()
	}else{
		nl.currentEl = nl.currentEl.Next()
	}
	return nl.currentEl.Value.(*Node)
}

// Prev returns the previous list node or nil.
func (nl *NodeList) Previous() *Node {
	if nl.currentEl == nil {
		nl.currentEl = nl.list.Front()
	}else{
		nl.currentEl = nl.currentEl.Prev()
	}
	return nl.currentEl.Value.(*Node)
}

// Back returns the last node of list or nil if the list is empty.
func (nl *NodeList) Back() *Node {
	return nl.list.Back().Value.(*Node)
}

// Front returns the first node of list or nil if the list is empty.
func (nl *NodeList) Front() *Node {
	return nl.list.Front().Value.(*Node)
}
