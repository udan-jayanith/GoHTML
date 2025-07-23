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

func WriteHTML(w io.Writer, node *Node) {
	rootNode := CreateNode(".")
	rootNode.Append(CreateNode("."))
	rootNode.AppendChild(node)

	stack := linkedliststack.New()
	traverser := GetTraverser(rootNode)

	traverser.Walkthrough(func(node *Node) {
		if node.GetTagName() == "" && node.GetNextNode() == nil{
			w.Write([]byte(node.GetText()))
			value, _ := stack.Pop()
			fmt.Fprintf(w, "</%s>", value.(string))
		} else if node.GetTagName() == "" {
			w.Write([]byte(node.GetText()))
		} else if node.GetChildNode() == nil && node.GetNextNode() == nil {
			fmt.Fprintf(w, "<%s></%s>", node.GetTagName(), node.GetTagName())
			value, _ := stack.Pop()
			fmt.Fprintf(w, "</%s>", value.(string))
			return
		} else if node.GetChildNode() == nil {
			fmt.Fprintf(w, "<%s></%s>", node.GetTagName(), node.GetTagName())
			return
		}else if node.GetChildNode() != nil {
			stack.Push(node.GetTagName())
			if node.GetTagName() != "."{
				fmt.Fprintf(w, "<%s>", node.GetTagName())
			}
		}
	})

}
