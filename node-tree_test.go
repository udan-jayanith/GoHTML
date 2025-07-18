package GoHtml_test

import (
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestNodeTree(t *testing.T){
	bodyEl := GoHtml.CreateEl("body")
	bodyEl.AppendChild(GoHtml.CreateEl("h1"))
	bodyEl.AppendChild(GoHtml.CreateEl("p"))

	traverser := GoHtml.GetTraverser(bodyEl.ChildNodes)
	for traverser.CurrentNode != nil{
		t.Log(traverser.CurrentNode.TagName)
		traverser.Next()
	}
}
