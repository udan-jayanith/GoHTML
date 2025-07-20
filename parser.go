package GoHtml

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

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
	traverser := GetTraverser(lastNode)
	for traverser.GetCurrentNode() != nil && traverser.GetCurrentNode().isClosed(){
		traverser.Previous()
	}
	return traverser.GetCurrentNode()
}
