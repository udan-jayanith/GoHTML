package Formatter

import (
	"github.com/rivo/uniseg"
)

func (fr *format_reader) wrap_text(text string) {
	if !fr.formatting_buf_is_last_write_eol() {
		fr.write_eol_to_formatting_buf()
	}

	b := []byte(text)
	state := -1
	var (
		c            []byte
		boundaries   int
		current_line []byte
	)

	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		current_line = append(current_line, c...)

		if boundaries&uniseg.MaskLine == uniseg.LineCanBreak &&
			fr.formatting_buf_line_length()+len(current_line) >= fr.format_options.LineLength {

			fr.write_string_to_formatting_buf(string(current_line))
			current_line = make([]byte, 0, len(current_line))
			fr.write_eol_to_formatting_buf()
		} else if boundaries&uniseg.MaskLine == uniseg.LineMustBreak {
			fr.write_string_to_formatting_buf(string(current_line))
			current_line = make([]byte, 0, len(current_line))
			fr.write_eol_to_formatting_buf()
		}
	}

	fr.write_string_to_formatting_buf(string(current_line))
	fr.write_eol_to_formatting_buf()
}
