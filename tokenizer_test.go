package GoHtml_test

import (
	"fmt"
	"net/http"

	GoHtml "github.com/udan-jayanith/GoHTML"
	"golang.org/x/net/html"
)

func ExampleTokenizer() {
	//Request the html
	res, err := http.Get("https://go.dev/")
	if err != nil || res.StatusCode != http.StatusOK {
		return
	}
	defer res.Body.Close()

	//NewTokenizer takes a io.reader that receives UTF-8 encoded html code and returns a Tokenizer.
	t := GoHtml.NewTokenizer(res.Body)
	//NewNodeTreeBuilder return a new NodeTreeBuilder that can be used to build a node tree.
	nodeTreeBuilder := GoHtml.NewNodeTreeBuilder()
	for {
		//Advanced scans the next token and returns its type.
		tt := t.Advanced()
		if tt == html.ErrorToken {
			break
		}

		//WriteNodeTree takes a node and a token type. The node can be nil so if token type is EndTagToken.
		nodeTreeBuilder.WriteNodeTree(t.GetCurrentNode(), tt)
	}

	//Prints the root node of the node tree in the nodeTreeBuilder.
	fmt.Println(nodeTreeBuilder.GetRootNode())
}
