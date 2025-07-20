package GoHtml_test

import (
	"testing"
	"os"
	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestDecode(t *testing.T){
	htmlFile, err := os.Open("./test-files/1.html")
	if err != nil {
		t.Fatal("Test file error: ", err)
		return
	} 
	
	GoHtml.Decode(htmlFile)
}