package GoHtml_test

import (
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestWalkthrough(t *testing.T) {
	body := GoHtml.CreateNode("body")
	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading")
	body.AppendChild(h1)
	p := GoHtml.CreateNode("p")
	p.AppendText("The HTML <p>tag is a fundamental element used for creating paragraphs in web development. It helps structure content, separating text into distinct blocks. When you wrap text within <p>... </p>tags, you tell browsers to treat the enclosed content as a paragraph.")
	body.AppendChild(p)


	traverser := GoHtml.NewTraverser(body)

	resList := make([]*GoHtml.Node, 0)
	traverser.Walkthrough(func(node *GoHtml.Node) GoHtml.TraverseCondition {
		resList = append(resList, node)
		return GoHtml.ContinueWalkthrough
	})

	testList := []*GoHtml.Node{
		body,
		h1, 
		h1.GetChildNode(),
		p,
		p.GetChildNode(),
	}
	for i := range testList {
		if testList[i] != resList[i]{
			t.Fatal("Expected ", testList[i], "but got ", resList[i], "in index ", i )
		}
	}
}
