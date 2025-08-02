package GoHtml

import (
	"strings"
	"sync"
)

type ClassList struct {
	classes map[string]struct{}
	rwMutex *sync.Mutex
}

// NewClassList returns a new empty ClassList.
func NewClassList() ClassList {
	classList := ClassList{
		classes: make(map[string]struct{}),
		rwMutex: &sync.Mutex{},
	}

	return classList
}

// AppendClass append className to classList. className that contains multiple classes is also a valid className.
func (classList ClassList) AppendClass(className string) {
	classList.rwMutex.Lock()
	defer classList.rwMutex.Unlock()

	classes := strings.SplitSeq(className, " ")
	for v := range classes {
		classList.classes[strings.TrimSpace(v)] = struct{}{}
	}
}

// SetClass append classes in the node to classList.
func (classList ClassList) SetClass(node *Node) {
	if node == nil {
		return 
	}
	classes, _ := node.GetAttribute("class")
	classList.AppendClass(classes)
}

// Contains returns whether the className exists or not.
func (classList ClassList) Contains(className string) bool {
	classList.rwMutex.Lock()
	defer classList.rwMutex.Unlock()

	classes := strings.SplitSeq(className, " ")
	for v := range classes {
		_, ok := classList.classes[strings.TrimSpace(v)]
		if !ok {
			return false
		}
	}

	return true
}

// DeleteClass deletes the specified classes in className.
func (classList ClassList) DeleteClass(className string) {
	classList.rwMutex.Lock()
	defer classList.rwMutex.Unlock()

	classes := strings.SplitSeq(className, " ")
	for v := range classes {
		delete(classList.classes, strings.TrimSpace(v))
	}
}

// Encode returns the full className.
func (classList ClassList) Encode() string {
	classList.rwMutex.Lock()
	defer classList.rwMutex.Unlock()

	classes := ""
	for v := range classList.classes {
		if classes != ""{
			classes+=" "
		}
		classes+=v
	}
	return classes
}

// EncodeTo encode className for the node.
func (classList ClassList) EncodeTo(node *Node){
	if node == nil {
		return
	}
	node.SetAttribute("class", classList.Encode())
}