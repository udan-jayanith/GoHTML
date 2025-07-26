package GoHtml_test

import (
	"os"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestEncodeToHTML(t *testing.T) {
	filePath := "./test-files/outputs/test-writeHTML.html"
	os.Remove(filePath)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
		return
	}

	body := GoHtml.CreateNode("body")
	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading")
	h1.SetAttribute("id", "first-heading'")
	body.AppendChild(h1)
	body.AppendChild(GoHtml.CreateNode("br"))
	p := GoHtml.CreateNode("p")
	p.AppendText("The HTML tag is a fundamental element used for creating paragraphs in web development. It helps structure content, separating text into distinct blocks. When you wrap text within tags, you tell browsers to treat the enclosed content as a paragraph.")
	body.AppendChild(p)

	GoHtml.EncodeToHTML(file, body)
}

func TestDecodeToNodeTree(t *testing.T) {
	file, err := os.Open("./test-files/1.html")
	if err != nil {
		t.Fatal(err)
		return
	}

	node, err := GoHtml.DecodeToNodeTree(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	filePath := "./test-files/outputs/TestDecodeToNodeTree.html"
	os.Remove(filePath)
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
		return
	}
	GoHtml.EncodeToHTML(file, node.GetNextNode())
}