package GoHtml_test

import(
	"testing"
	"github.com/udan-jayanith/GoHTML"
	"net/http"
	"time"
)
/*
Adapted from [GoQuery example](https://github.com/PuerkitoBio/goquery?tab=readme-ov-file#examples)
*/
func TestFetchPostCovers(t *testing.T){
	res, err := http.Get("https://www.metalsucks.net/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	tim := time.Now()
	node, err := GoHtml.Decode(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	nodeList := node.QuerySelectorAll(".left-content article .post-title")
	t.Log("Got ", nodeList.Len(), " post titles.")
	iter := nodeList.IterNodeList()
	for node := range iter{
		t.Log("---------Post title-----------")
		t.Log(node.GetInnerText())
	}
	t.Log(time.Since(tim).Seconds())
}

func toNodeTree(url string) *GoHtml.Node{
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK{
		return nil
	}
	defer res.Body.Close()

	rootNode, _ := GoHtml.Decode(res.Body)
	return rootNode
}