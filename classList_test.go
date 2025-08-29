package GoHtml_test

import (
	"fmt"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestClasses(t *testing.T) {
	node := GoHtml.CreateNode("div")
	node.SetAttribute("class", "div-container main")

	classList := GoHtml.NewClassList()
	classList.DecodeFrom(node)
	if !classList.Contains("main") {
		t.Fatal("")
		return
	}
	classList.DeleteClass("main")
	if classList.Contains("main") {
		t.Fatal("")
		return
	}

	classList.AppendClass("main-div")
	if !classList.Contains("main-div") {
		t.Fatal("")
		return
	}

	classList.EncodeTo(node)
}

func ExampleClassList_Contains() {
	//Creates a div that has classes video-container and main-contents
	div := GoHtml.CreateNode("div")
	div.SetAttribute("class", "video-container main-contents")

	classList := GoHtml.NewClassList()
	//Add the classes in the div to the class list
	classList.DecodeFrom(div)

	//Checks wether the following classes exists in the classList
	fmt.Println(classList.Contains("container"))
	fmt.Println(classList.Contains("video-container"))

	//Output:
	//false
	//true
}

func ExampleClassList_Encode(){
	classList := GoHtml.NewClassList()

	//Add classes to the class list
	classList.AppendClass("container")
	classList.AppendClass("warper")
	classList.AppendClass("main-content")

	//This would output something like this "warper container main-content". Order of the output is not guaranteed.
	fmt.Println(classList.Encode())
}