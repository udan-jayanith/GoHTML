package Formatter_test

import(
	"testing"
	"os"
	"github.com/udan-jayanith/GoHTML/formatter"
	"io"
)

func TestFormat(t *testing.T){
	file, err := os.Open("../test-files/2.html") 
	if err != nil {
		t.Fatal(err.Error())
	}
	
	r, err := Formatter.Format(file)
	byts, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err.Error())
	}
	
	t.Log(string(byts), len(byts))
}

