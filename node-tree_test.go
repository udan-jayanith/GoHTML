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

	traverser := GoHtml.GetTraverser(bodyEl.GetChildNode())
	for traverser.GetCurrentNode() != nil{
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
	
	for traverser.GetCurrentNode() != nil{
		if traverser.GetCurrentNode().GetTagName() != fmt.Sprintf("h%v", count){
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

	if body.GetChildNode().GetLastNode().GetTagName() != "footer"{
		t.FailNow()
	}
}

func TestGetFirstNode(t *testing.T){
	body := GoHtml.CreateNode("body")
	body.AppendChild(GoHtml.CreateNode("h1"))
	body.AppendChild(GoHtml.CreateNode("p"))
	body.AppendChild(GoHtml.CreateNode("footer"))

	if body.GetChildNode().GetFirstNode().GetTagName() != "h1"{
		t.FailNow()
	}
}

func TestAppendTextAndInnerText(t *testing.T){
	p := GoHtml.CreateNode("body")
	text := "Hello world"
	p.AppendText(text)

	if p.GetInnerText() != text{
		t.FailNow()
	}
}