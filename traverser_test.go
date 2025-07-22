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

	traverser := GoHtml.GetTraverser(body)
	traverser.Walkthrough(func(node *GoHtml.Node) {
		t.Log(node)
	})
}
