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

// This will add appropriate spacing to behind of attributes.
func format_attribute_list(iterator iter.Seq[string], line_length int, buf formatting_buf_writers) {
	for attribute := range iterator {
		if buf.formatting_buf_is_last_write_eol() {
			buf.write_string_to_formatting_buf(attribute)
		} else if buf.formatting_buf_line_length()+len(attribute) >= line_length {
			buf.write_eol_to_formatting_buf()
			buf.write_string_to_formatting_buf(attribute)
		} else {
			buf.write_string_to_formatting_buf(" " + attribute)
		}
	}
}