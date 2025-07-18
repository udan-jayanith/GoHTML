package GoHtml

func CreateEl(tagName string) *Node {
	return &Node{
		TagName: tagName,
		Closed:  true,
	}
}
