package GoHtml_test

import (
	"os"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestIterNodeList(t *testing.T) {
	file, err := os.Open("./test-files/3.html")
	if err != nil {
		t.Fatal(err)
		return
	}

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	list := GoHtml.NewNodeList()
	traverser := GoHtml.GetTraverser(node)
	traverser.Walkthrough(func(node *GoHtml.Node) GoHtml.TraverseCondition {
		list.Append(node)
		return GoHtml.ContinueWalkthrough
	})

	iterator := list.IterNodeList()
	for node := range iterator{
		t.Log(node.GetTagName())
	}
}
