package GoHtml_test

import (
	"os"
	"strings"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

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
