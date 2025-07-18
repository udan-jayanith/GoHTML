package GoHtml

/*
A void element is an element in HTML that cannot have any child nodes (i.e., nested elements or text nodes). Void elements only have a start tag; end tags must not be specified for void elements.
*/

import (
	"strings"
)

const (
	Area          = "area"
	Base   string = "base"
	Br     string = "br"
	Col    string = "col"
	Embed  string = "embed"
	Hr     string = "hr"
	Img    string = "img"
	Input  string = "input"
	Link   string = "link"
	Meta   string = "meta"
	Param  string = "param"
	Source string = "source"
	Track  string = "track"
	Wbr    string = "wbr"
)

var (
	VoidTags = map[string]bool{
		"area":   true,
		"base":   true,
		"br":     true,
		"col":    true,
		"embed":  true,
		"hr":     true,
		"img":    true,
		"input":  true,
		"link":   true,
		"meta":   true,
		"param":  true,
		"source": true,
		"track":  true,
		"wbr":    true,
		"!doctype": true,
	}
)

//IsVoidTag returns whether the tagName is a void tag or DTD 
func IsVoidTag(tagName string) bool{
	tagName = strings.TrimSpace(strings.ToLower(tagName))
	return VoidTags[tagName]
}

/*
 A DTD defines the structure and the legal elements and attributes of an XML document.
*/
const (
	//This is not a void el. but added it anyway.
	DOCTYPEDTD string = "!DOCTYPE"
)
