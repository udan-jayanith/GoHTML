package GoHtml_test

import (
	"os"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestGetElementBy(t *testing.T) {
	file, err := os.Open("./test-files/3.html")
	if err != nil {
		t.Fatal(err)
		return
	}

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	title := node.GetElementByTagName("title")
	if title == nil {
		t.Fatal("title is nil.")
		return
	} else if title.GetInnerText() != "Simple Page" {
		t.Fatal("Unexpected title.")
		return
	}

	heading := node.GetElementByClassName("heading")
	if heading == nil {
		t.Fatal("heading is nil.")
		return
	} else if heading.GetTagName() != "h1" {
		t.Fatal("Exacted tag name is p but got ", heading.GetTagName())
		return
	}
	classList := GoHtml.NewClassList()
	classList.SetClass(heading)

	if !classList.Contains("heading") {
		t.Fatal("Exacted class name heading but go ", classList.Encode())
	}

	paragraph := node.GetElementByClassName("paragraph")
	if paragraph != nil {
		t.Fatal("Got none existing node.", paragraph)
		return
	}
}

func TestGetElementByID(t *testing.T){
	file, err := os.Open("./test-files/3.html")
	if err != nil {
		t.Fatal(err)
		return
	}

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	heading := node.GetElementByID("hello-world")
	if heading == nil {
		t.Fatal("Heading is nil")
		return
	}
	id, _ := heading.GetAttribute("id")
	if id != "hello-world"{
		t.Fatal("unexpected id", id, heading)
		return
	}
}