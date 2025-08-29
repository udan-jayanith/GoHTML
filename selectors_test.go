package GoHtml_test

import (
	"testing"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

func TestTokenizeSelector(t *testing.T) {
	slice := GoHtml.TokenizeSelectorsAndCombinators(".class-1 > .class-2 + .class-3 a")
	if len(slice) != 4 {
		t.Fatal("Exacted slice length of", 4, "but got", len(slice))
	}
}
