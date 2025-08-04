package GoHtml_test

import (
	"os"
	"testing"

	"github.com/emirpasic/gods/stacks/linkedliststack"
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
	classList.DecodeFrom(heading)

	if !classList.Contains("heading") {
		t.Fatal("Exacted class name heading but go ", classList.Encode())
	}

	paragraph := node.GetElementByClassName("paragraph")
	if paragraph != nil {
		t.Fatal("Got none existing node.", paragraph)
		return
	}
}

func TestGetElementByID(t *testing.T) {
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
	if id != "hello-world" {
		t.Fatal("unexpected id", id, heading)
		return
	}
}

func testFile4NodeTree() (*GoHtml.Node, error) {
	file, err := os.Open("./test-files/4.html")
	if err != nil {
		return nil, err
	}

	node, err := GoHtml.Decode(file)
	return node, err
}

func TestGetElementsByClassName(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
		return
	}

	nodeList := node.GetElementsByClassName("ordered-item")
	iterator := nodeList.IterNodeList()
	stack := linkedliststack.New()
	stack.Push("Mango")
	stack.Push("Orange")
	stack.Push("Apple")

	if nodeList.Len() == 0 {
		t.Fatal("NodeList is empty")
	}

	for node := range iterator {
		value, ok := stack.Pop()
		if !ok {
			t.Fatal("Stack is empty")
		}
		text := value.(string)
		if node.GetInnerText() != text {
			t.Fatal("Unexpected text", node)
		}
	}
}

func TestGetElementsByTagName(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
		return
	}

	nodeList := node.GetElementsByClassName("ordered-item")
	iterator := nodeList.IterNodeList()
	stack := linkedliststack.New()
	stack.Push("Mango")
	stack.Push("Orange")
	stack.Push("Apple")

	if nodeList.Len() == 0 {
		t.Fatal("NodeList is empty")
	}

	for node := range iterator {
		value, ok := stack.Pop()
		if !ok {
			t.Fatal("Stack is empty")
		}
		text := value.(string)
		if node.GetInnerText() != text {
			t.Fatal("Unexpected text", node)
		}
	}
}

func TestGetElementsById(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
		return
	}

	nodeList := node.GetElementsByClassName("ordered-item")
	iterator := nodeList.IterNodeList()
	stack := linkedliststack.New()
	stack.Push("Mango")
	stack.Push("Orange")
	stack.Push("Apple")

	if nodeList.Len() == 0 {
		t.Fatal("NodeList is empty")
	}

	for node := range iterator {
		value, ok := stack.Pop()
		if !ok {
			t.Fatal("Stack is empty")
		}
		text := value.(string)
		if node.GetInnerText() != text {
			t.Fatal("Unexpected text", node)
		}
	}
}
