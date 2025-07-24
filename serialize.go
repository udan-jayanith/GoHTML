package GoHtml

import (
	//"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/stacks/linkedliststack"
)

/*
func Decode(rd io.Reader){
	reader := bufio.NewReader(rd)
	var currentNode *Node

	strBuf := ""
	quote := ""
	for {
		byt, err := reader.ReadByte()
		if err != nil{
			break
		}

		strBuf += string(byt)

		if IsQuote(string(byt)) && quote == "" && !regexp.MustCompile(`^\s*<.+$`).MatchString(strBuf) {
			quote = string(byt)
		}else if IsQuote(string(byt)) && quote == string(byt){
			quote = ""
		}

		if quote != ""{
			continue
		}


		htmlTagRegex := regexp.MustCompile(`^\s*<.+>$`)
		htmlTextRegex := regexp.MustCompile(`<$`)

		if htmlTagRegex.MatchString(strBuf) {
			strBuf = ""
		}else if htmlTextRegex.MatchString(strBuf) && len(strings.TrimSpace(strBuf)) > 1{
			strBuf = strBuf[:len(strBuf)-1]
			if currentNode == nil{
				currentNode = CreateNode("")
				currentNode.SetText(strBuf)
			}else{
				currentNode.AppendText(strBuf)
			}

			strBuf = "<"
		}
	}
}

func SerializeHTML(tag string){
	//This is the regex that will be used to parse the tag
	//(\w+(?:-\w+)*)\s*(?:=\s*(?:(["'`+"`"+`])(.*?)\2|(\S+)))?
}

func IsClosingTag(tag string) bool {
	reg := regexp.MustCompile(`^<\/.*>\s*$`)
	return reg.MatchString(tag)
}

func IsQuote(chr string) bool {
	return chr == `"` || chr == `'` || chr == "`"
}

func getFirstUnclosedTagInNodeChain(lastNode *Node) *Node{
}

*/

func wrapAttributeValue(value string) string {
	reg := regexp.MustCompile(`^[\d\.]+$`)
	if reg.Match([]byte(value)) {
		return value
	}

	return `"` + strings.ReplaceAll(value, `"`, "&quot;") + `"`
}

func decodeListAttributes(node *Node) string {
	w := strings.Builder{}
	node.IterateAttributes(func(attribute, value string) {
		w.Write([]byte(fmt.Sprintf(" %s=%s", attribute, wrapAttributeValue(value))))
	})
	return w.String()
}

func WriteHTML(w io.Writer, rootNode *Node) {
	type stackFrame struct {
		node      *Node
		openedTag bool
	}

	stack := linkedliststack.New()
	stack.Push(stackFrame{node: rootNode, openedTag: false})

	for stack.Size() > 0 {
		t, _ := stack.Pop()
		top := t.(stackFrame)
		current := top.node

		if current == nil {
			continue
		}

		tagName := current.GetTagName()
		if tagName == "" {
			w.Write([]byte(current.GetText()))
		} else if IsVoidTag(tagName) {
			fmt.Fprintf(w, "<%s %s>", tagName, decodeListAttributes(current))
			if current.GetNextNode() != nil {
				stack.Push(stackFrame{node: current.GetNextNode(), openedTag: false})
			}
		} else if !top.openedTag {
			fmt.Fprintf(w, "<%s %s>", tagName, decodeListAttributes(current))
			stack.Push(stackFrame{node: current, openedTag: true})

			if current.GetChildNode() != nil {
				stack.Push(stackFrame{node: current.GetChildNode(), openedTag: false})
			}
		} else {
			fmt.Fprintf(w, "</%s>", tagName)
			if current.GetNextNode() != nil {
				stack.Push(stackFrame{node: current.GetNextNode(), openedTag: false})
			}
		}
	}
}
