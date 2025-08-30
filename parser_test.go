package GoHtml_test

import (
	"fmt"
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
	defer file.Close()

	node, err := GoHtml.Decode(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	var builder strings.Builder
	GoHtml.Encode(&builder, node)
}

func ExampleDecode() {
	r := strings.NewReader(`
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>User Profile</title>
		</head>
		<body>
			<h1 class="username">Udan</h1>
			<p class="email">udanjayanith@gmail.com</p>
			<p>Joined: 01/08/2024</p>
		</body>
	</html>
	`)

	rootNode, _ := GoHtml.Decode(r)

	titleNode := rootNode.QuerySelector("title")
	title := ""
	if titleNode != nil {
		title = titleNode.GetInnerText()
	}
	fmt.Println(title)
	//Output: 
	//User Profile
}
