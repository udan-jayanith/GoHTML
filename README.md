# GoHTML

A powerful and comprehensive HTML parser and DOM manipulation library for Go, bringing JavaScript-like DOM operations to the Go ecosystem.
Note: GoHTML only support UTF-8. It's users responsibility to make sure input is UTF-8.

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

	//Parses the given html reader and then returns the root node and an error.
	node, err := GoHtml.Decode(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	nodeList := node.GetElementsByClassName("post-title")
	iter := nodeList.IterNodeList()
	for node := range iter{
		print(node.GetInnerText())
	}
```

## Changelog

Changes, bug fixes and new features in this version.
- add: Tokenizer
- add: NodeTreeBuilder
- renamed: QuerySelector to Query
- renamed: QuerySelectorAll to QueryAll

## Documentation

Fully fledged [documentation](https://pkg.go.dev/github.com/udan-jayanith/GoHTML) is available at [go.pkg](https://pkg.go.dev/)

## Contributions

Contributions are welcome and pull requests and issues will be viewed by an official.
