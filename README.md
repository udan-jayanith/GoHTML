# [GoHTML](https://pkg.go.dev/github.com/udan-jayanith/GoHTML)

A HTML parse and a serializer for Go. GoHTML tries to keep semantic similar to JS-DOM API while trying to keep the API simple by not forcing JS-DOM model into GoHTML. Because of this GoHTML has node tree model. GoHTML tokenizer uses std net/html module for tokenizing in underlining layer. There for it's users responsibility to make sure inputs to GoHTML is UTF-8 encoded. GoHTML allows direct access to the node tree.

Note: Under development.
## Installation

Run the following command in project directory in order to install.

```bash
go get github.com/udan-jayanith/GoHTML
```

Then GoHTML can import like this.

```go
import (
	GoHtml "github.com/udan-jayanith/GoHTML"
)
```

## Features

- Parsing
- Serialization
- Node tree traversing
- Querying

## Example

Heres an example of fetching a website and parsing and then using querying methods.

```go
	res, err := http.Get("https://www.metalsucks.net/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	node, err := GoHtml.Decode(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	nodeList := node.QuerySelectorAll(".left-content article .post-title")
	for node := range nodeList.IterNodeList(){
		println(node.GetInnerText())
	}
```

## Documentation

Fully fledged [documentation](https://pkg.go.dev/github.com/udan-jayanith/GoHTML) is available at [go.pkg](https://pkg.go.dev/)

## Contributions

Contributions are welcome and pull requests and issues will be viewed by an official.
