package Formatter

import (
	"fmt"
	"iter"
	"strings"

	"golang.org/x/net/html"
)

func format_kv(key, value string) (str string) {
	if strings.TrimSpace(value) == "" {
		str = value
	} else {
		str = fmt.Sprintf(`%s="%s"`, key, value)
	}
	return str
}

func tokenize_kv_map(iterator map[string]string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k, v := range iterator {
			if !yield(format_kv(k, v)) {
				break
			}
		}
	}
}

func tokenize_kv_attr(attr []html.Attribute) iter.Seq[string]{
	return func(yield func(string) bool) {
		for _, pair := range attr {
			k, v := pair.Key, pair.Val
			if !yield(format_kv(k, v)) {
				break
			}
		}
	}
}

func format_attribute_list(iterator iter.Seq[string], byte_per_line uint, end EOL) string {
	var str string
	
}