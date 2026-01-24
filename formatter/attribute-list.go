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

func tokenize_kv_node(iterator map[string]string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k, v := range iterator {
			if !yield(format_kv(k, v)) {
				break
			}
		}
	}
}

func tokenize_kv_attr(attr []html.Attribute) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, pair := range attr {
			k, v := pair.Key, pair.Val
			if !yield(format_kv(k, v)) {
				break
			}
		}
	}
}

func format_attribute_list(iterator iter.Seq[string], bytes_per_line int, buf formatting_buf_writers) {
	var line_length int = 0
	for attribute := range iterator {
		if line_length >= bytes_per_line {
			buf.write_eol_to_formatting_buf()
			line_length = 0
		}

		if line_length == 0 {
			buf.write_string_to_formatting_buf(attribute)
		} else {
			buf.write_string_to_formatting_buf(" " + attribute)
			line_length++
		}
		line_length += len(attribute)
	}
}
