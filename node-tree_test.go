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

	traverser := GoHtml.NewTraverser(bodyEl.GetChildNode())
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
	traverser := GoHtml.NewTraverser(h1)
	
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
		t.Fatal("Unexpected parent node", hr.GetParent())
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
	body := GoHtml.CreateNode("body")
	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading")
	body.AppendChild(h1)
	p := GoHtml.CreateNode("p")
	p.AppendText("The HTML <p>tag is a fundamental element used for creating paragraphs in web development. It helps structure content, separating text into distinct blocks. When you wrap text within <p>... </p>tags, you tell browsers to treat the enclosed content as a paragraph.")
	body.AppendChild(p)

	if body.GetInnerText() != h1.GetChildNode().GetText() + p.GetChildNode().GetText(){
		t.Fatal(body.GetInnerText(), " != ", h1.GetChildNode().GetText() + p.GetChildNode().GetText())
	}

}

func TestRemoveNode(t *testing.T){
	article := GoHtml.CreateNode("article")
	
	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading.")
	article.AppendChild(h1)

	article.AppendChild(GoHtml.CreateNode(GoHtml.Br))

	p := GoHtml.CreateNode("p")
	p.AppendText("this is a paragraph.")
	article.AppendChild(p)

	h1.RemoveNode()

	if article.GetChildNode().GetTagName() != GoHtml.Br{
		t.Fatal("Unexpected tag. ", article.GetChildNode().GetTagName())
		return
	}else if p.GetParent() != article {
		t.Fatal("Unexpected parent.")
	}

	//p.RemoveNode()
	//t.Log(GoHtml.NodeTreeToHTML(article))
}

/*
func TestClosest(t *testing.T){
	node, err := testFile4NodeTree()
	if err != nil{
		t.Fatal(err)
	}
	node = node.GetElementByClassName("ordered-item")
	if node == nil {
		t.Fatal("Node is nil.")
	}

	node = node.Closest("ol .ordered-list")
	if node == nil {
		t.Fatal("Node is nil")
	}else if node.GetTagName() != "ol"{
		t.Fatal("Unexpected element.")
	}


}
*/