![GoHTML logo](https://raw.githubusercontent.com/udan-jayanith/GoHTML/46044619ab943b8ae00301565cc37566d5f2ffa4/assets/media/Black-text%20version.svg)
# GoHTML
A powerful and comprehensive HTML parser and DOM manipulation library for Go, bringing JavaScript-like DOM operations to the Go ecosystem.

## Installation
Run the following command in project directory in order to install. 
```bash
go get github.com/udan-jayanith/GoHTML
```

Then import like this
```go
import (
	GoHtml "github.com/udan-jayanith/GoHTML"
)
```

## Features
 * Parsing
 * Serialization
 * Node tree tarversing
 * Querying
 * [Concurrency safety](#concurrency-safety)

## Concurrency safety
Node tree itself is not concurrency safe but nodes and everything other then node trees are concurrency safe.

## Documentation
Fully fledged [documentation](https://pkg.go.dev/github.com/udan-jayanith/GoHTML) is available at [go.pkg](https://pkg.go.dev/)

## Contributions
Contributions are welcome and pull requests and issues will be viewed by an official.
