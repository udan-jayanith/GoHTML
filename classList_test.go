package GoHtml_test

import(
	"testing"
	"github.com/udan-jayanith/GoHTML"
)

func TestClasses(t *testing.T){
	node := GoHtml.CreateNode("div")
	node.SetAttribute("class", "div-container main")

	classList := GoHtml.NewClassList()
	classList.DecodeFrom(node)
	if !classList.Contains("main"){
		t.Fatal("")
		return
	}
	classList.DeleteClass("main")
	if classList.Contains("main"){
		t.Fatal("")
		return
	}

	classList.AppendClass("main-div")
	if !classList.Contains("main-div"){
		t.Fatal("")
		return
	}

	classList.EncodeTo(node)
}