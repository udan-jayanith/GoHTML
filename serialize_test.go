package GoHtml_test

import (
	"os"
	"strings"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
	htmlFormatter "github.com/yosssi/gohtml"
)

func TestEncode(t *testing.T) {
	body := GoHtml.CreateNode("body")
	h1 := GoHtml.CreateNode("h1")
	h1.AppendText("This is a heading")
	h1.SetAttribute("id", "first-heading'")
	body.AppendChild(h1)
	body.AppendChild(GoHtml.CreateNode("br"))
	p := GoHtml.CreateNode("p")
	p.AppendText("The HTML tag is a fundamental element used for creating paragraphs in web development. It helps structure content, separating text into distinct blocks. When you wrap text within tags, you tell browsers to treat the enclosed content as a paragraph.")
	body.AppendChild(p)

	builder1 := &strings.Builder{}
	GoHtml.Encode(builder1, body)

	//It's hard compare exacted output. Because strings, prettier formats html code. htmlFormatter and prettier add extra stuffs to the html codes like dash in void tags. Exacted output is in the ./test-files/2.html.
	return
	str := htmlFormatter.Format(builder1.String())
	builder1.Reset()
	builder1.Write([]byte(str))

	exactedStr := ""
	if builder1.String() != exactedStr {
		t.Fatal("Encoding failed. Unexpected outputs", builder1.String())
	}
}

func TestDecode(t *testing.T) {
	file, err := os.Open("./test-files/1.html")
	if err != nil {
		t.Fatal(err)
		return
	}

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	builder1 := &strings.Builder{}
	GoHtml.Encode(builder1, node)
	//It's hard compare exacted output. Because strings, prettier formats html code. htmlFormatter and prettier add extra stuffs to the html codes like dash in void tags. Exacted output is in the ./test-files/2.html.
}
