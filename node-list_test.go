package GoHtml_test

import (
	"fmt"
	"os"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestIterNodeList1(t *testing.T) {
	file, err := os.Open("./test-files/3.html")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer file.Close()

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	list := GoHtml.NewNodeList()
	traverser := GoHtml.NewTraverser(node)
	traverser.Walkthrough(func(node *GoHtml.Node) GoHtml.TraverseCondition {
		list.Append(node)
		return GoHtml.ContinueWalkthrough
	})

	iterator := list.IterNodeList()
	for node := range iterator{
		node.RemoveNode()
	}
}

func TestIterNodeList2(t *testing.T){
	nodeList := GoHtml.NewNodeList()
	iter := nodeList.IterNodeList()
	for node := range iter{
		t.Log(node)
	}
}

func ExampleNodeList(){
	nodeList := GoHtml.NewNodeList()
	nodeList.Append(GoHtml.CreateNode("br"))
	nodeList.Append(GoHtml.CreateNode("hr"))
	nodeList.Append(GoHtml.CreateNode("div"))

	iter := nodeList.IterNodeList()
	for node := range iter{
		fmt.Println(node.GetTagName())
	}
	//Output: 
	//br
	//hr 
	//div
}