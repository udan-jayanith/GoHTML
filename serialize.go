package GoHtml

import (
	//"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
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

type visitedNode struct {
	childNodes bool
	nextNodes  bool
}

func WriteHTML(w io.Writer, rootNode *Node) {
	traverser := GetTraverser(rootNode)
	for traverser.GetCurrentNode() != nil {
		if traverser.GetCurrentNode().GetTagName() == "" {
			w.Write([]byte(traverser.GetCurrentNode().GetText()))
		} else {
			fmt.Fprintf(w, "<%s>", traverser.GetCurrentNode().GetTagName())
			if traverser.GetCurrentNode().GetChildNode() != nil {
				WriteHTML(w, traverser.GetCurrentNode().GetChildNode())
			}
			fmt.Fprintf(w, "</%s>", traverser.GetCurrentNode().GetTagName())
		}

		traverser.Next()
	}
}
