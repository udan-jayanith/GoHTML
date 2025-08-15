package GoHtml_test

import(
	"testing"
	"github.com/udan-jayanith/GoHTML"
	"net/http"
	"time"
)

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

	nodeList := node.QueryAll(".sm-feat .clearfix article")
	t.Log("Got ", nodeList.Len(), " post titles.")
	iter := nodeList.IterNodeList()
	for node := range iter{
		t.Log("---------Post title-----------")
		t.Log(node.GetInnerText())
	}
	t.Log(time.Since(tim).Seconds())
}