package GoHtml_test

import (
	"fmt"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestNodeTree(t *testing.T){
	bodyEl := GoHtml.CreateNode("body")
	bodyEl.AppendChild(GoHtml.CreateNode("h1"))
	bodyEl.AppendChild(GoHtml.CreateNode("p"))

	traverser := GoHtml.GetTraverser(bodyEl.ChildNode)
	for traverser.CurrentNode != nil{
		traverser.Next()
	}
}

func TestAppend(t *testing.T){
	h1 := GoHtml.CreateNode("h1")
	for i:=2; i<=6; i++{
		h1.Append(GoHtml.CreateNode(fmt.Sprintf("h%v", i)))
	}

	count := 1
	traverser := GoHtml.GetTraverser(h1)
	for traverser.CurrentNode != nil{
		if traverser.CurrentNode.TagName != fmt.Sprintf("h%v", count){
			t.Fatal("Unexpected tag name.")
		}
		
		traverser.Next()
		count++
	}
}

func TestGetParent(t *testing.T){
	article := GoHtml.CreateNode("article")
	article.AppendChild(GoHtml.CreateNode("h1"))

	p := GoHtml.CreateNode("p")
	article.AppendChild(p)
	hr := GoHtml.CreateNode("hr")
	article.AppendChild(hr)

	if hr.GetParent() != article || p.GetParent() != article{
		t.Fatal("Unexpected parent node")
	}
}

func TestGetLastNode(t *testing.T){
	body := GoHtml.CreateNode("body")
	body.AppendChild(GoHtml.CreateNode("h1"))
	body.AppendChild(GoHtml.CreateNode("p"))
	body.AppendChild(GoHtml.CreateNode("footer"))

	if body.ChildNode.GetLastNode().TagName != "footer"{
		t.FailNow()
	}
}

func TestGetFirstNode(t *testing.T){
	body := GoHtml.CreateNode("body")
	body.AppendChild(GoHtml.CreateNode("h1"))
	body.AppendChild(GoHtml.CreateNode("p"))
	body.AppendChild(GoHtml.CreateNode("footer"))

	if body.ChildNode.GetFirstNode().TagName != "h1"{
		t.FailNow()
	}
}

func TestAppendTextAndInnerText(t *testing.T){
	body := GoHtml.CreateNode("body")

	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading")
	body.AppendChild(h1)

	p := GoHtml.CreateNode("p")
	p.AppendText(" Hello world")
	body.AppendChild(p)

	body.AppendChild(GoHtml.CreateNode("footer"))

	if body.GetInnerText() != "This is a heading Hello world"{
		t.FailNow()
	}
}