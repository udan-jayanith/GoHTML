package GoHtml_test

import(
	"testing"
	"github.com/udan-jayanith/GoHTML"
)

func TestTokenizeSelector(t *testing.T){
	slice := GoHtml.TokenizeSelectorsAndCombinators(".class-1 > .class-2 + .class-3 a")
	for _, el := range slice{
		t.Log(el)
	}
}