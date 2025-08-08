package GoHtml_test

import (
	"os"
	"strings"
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestDecode(t *testing.T) {
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

	var builder strings.Builder
	GoHtml.Encode(&builder, node)
	t.Log(builder.String())
}
