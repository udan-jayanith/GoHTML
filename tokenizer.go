package GoHtml
import (
	"bufio"
	"io"

	"github.com/emirpasic/gods/stacks/linkedliststack"
	"golang.org/x/net/html"
)

type TokenType = html.TokenType

const (
	// ErrorToken means that an error occurred during tokenization.
	ErrorToken TokenType = html.ErrorToken
	// TextToken means a text node.
	TextToken = html.TextToken
	// A StartTagToken looks like <a>.
	StartTagToken = html.StartTagToken
	// An EndTagToken looks like </a>.
	EndTagToken = html.EndTagToken
	// A SelfClosingTagToken tag looks like <br/>.
	SelfClosingTagToken = html.SelfClosingTagToken
	// A CommentToken looks like <!--x-->.
	CommentToken = html.CommentToken
	// A DoctypeToken looks like <!DOCTYPE x>
	DoctypeToken = html.DoctypeToken
)

type Attribute struct {
	Key, Val string
}

type Token struct {
	Type TokenType
	Data string
	Attr []Attribute
}

// Tokenizer contains a *html.Tokenizer.
type Tokenizer struct {
	rd      *bufio.Reader
	buf     []byte
	readBuf []byte
}

// NewTokenizer returns a new Tokenizer.
func NewTokenizer(r io.Reader) Tokenizer {
	return Tokenizer{
		rd:      bufio.NewReader(r),
		buf:     make([]byte, 0, 32),
		readBuf: make([]byte, 0, 32),
	}
}

func (t *Tokenizer) Token() Token {
	panic("Not implemented")
	return Token{}
}

// Advanced scans the next token and returns its type.
func (t *Tokenizer) Advanced() TokenType {
	//return t.z.Next()
	panic("Not implemented")
	return ErrorToken
}

// Tags, closingTags and comments are enclosed by < >.
func (t *Tokenizer) scan() {
	buffer := make([]byte, 1024)
	for {
		buf := make([]byte, 1024)
		n, err := t.rd.Read(buf)
		//handle the error
		if err != nil {
			break
		}

		for i := 0; i < n; i++ {
			byt := t.bytesBuf[i]

		}
	}
}


// NodeTreeBuilder is used to build a node tree given a node and it's type.
type NodeTreeBuilder struct {
	rootNode    *Node
	stack       *linkedliststack.Stack
	currentNode *Node
}

// NewNodeTreeBuilder returns a new NodeTreeBuilder.
func NewNodeTreeBuilder() NodeTreeBuilder {
	rootNode := CreateTextNode("")
	return NodeTreeBuilder{
		rootNode:    rootNode,
		currentNode: rootNode,
		stack:       linkedliststack.New(),
	}
}

// WriteNodeTree append the node given html.TokenType.
func (ntb *NodeTreeBuilder) WriteNodeTree(node *Node, tt html.TokenType) {
	switch tt {
	case html.EndTagToken:
		val, ok := ntb.stack.Pop()
		if !ok || val == nil {
			return
		}
		ntb.currentNode = val.(*Node)
	case html.DoctypeToken, html.StartTagToken, html.SelfClosingTagToken, html.TextToken:
		if node == nil {
			return
		}

		if isTopNode(ntb.currentNode, ntb.stack) {
			ntb.currentNode.AppendChild(node)
		} else {
			ntb.currentNode.Append(node)
		}

		if !node.IsTextNode() && !IsVoidTag(node.GetTagName()) {
			ntb.stack.Push(node)
		}
		ntb.currentNode = node
	}
}

// GetRootNode returns the root node of the accumulated node tree and resets the NodeTreeBuilder.
func (ntb *NodeTreeBuilder) GetRootNode() *Node {
	node := ntb.rootNode.GetNextNode()
	ntb.rootNode.RemoveNode()

	rootNode := CreateTextNode("")
	ntb.rootNode = rootNode
	ntb.currentNode = rootNode
	ntb.stack = linkedliststack.New()

	return node
}

func isTopNode(node *Node, stack *linkedliststack.Stack) bool {
	val, ok := stack.Peek()
	if !ok || val == nil {
		return false
	}

	topNode := val.(*Node)
	return topNode == node
}