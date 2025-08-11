package GoHtml_test

import (
	"os"
	"testing"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	GoHtml "github.com/udan-jayanith/GoHTML"
)

func testFile4NodeTree() (*GoHtml.Node, error) {
	file, err := os.Open("./test-files/4.html")
	if err != nil {
		return nil, err
	}

	node, err := GoHtml.Decode(file)
	return node, err
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

func TestGetElementByClassName(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
	}

	node = node.GetElementByClassName("ordered-item")
	if node == nil {
		t.Fatal("Node is nil")
	}else if node.GetInnerText() != "Apple" {
		t.Fatal("Expected Apple but got ", node.GetInnerText())
	}
}

func TestGetElementByTagName(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
	}

	node = node.GetElementByTagName("h2")
	if node == nil {
		t.Fatal("Node is nil")
	}else if node.GetInnerText() != "List 1"{
		t.Fatal("Expected List 1 but got ", node.GetInnerText())
	}
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

	for node := range iterator {
		value, ok := stack.Pop()
		if !ok {
			t.Fatal("Stack is empty")
		}
		text := value.(string)
		if node.GetInnerText() != text {
			t.Fatal("Expected ", text, " But got ", node.GetInnerText())
		}
	}
	if nodeList.Len() == 0 {
		t.Fatal("NodeList is empty")
	}
}

func TestGetElementsByTagName(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
		return
	}

	nodeList := node.GetElementsByTagName("meta")
	if nodeList.Len() != 2 {
		t.Fatal(nodeList.Len())
	}
}

func TestGetElementsById(t *testing.T) {
	node, err := testFile4NodeTree()
	if err != nil {
		t.Fatal(err)
		return
	}

	nodeList := node.GetElementsById("idElement")
	iter := nodeList.IterNodeList()
	stack := linkedliststack.New()
	stack.Push("Lorem")
	stack.Push("")

	for node := range iter {
		val, ok := stack.Pop()
		if !ok {
			t.Fatal("Stack error.")
		}

		if node.GetInnerText() != val.(string) {
			t.Fatal("Unexpected node: ", node.GetInnerText(), val.(string))
		}
	}
}
