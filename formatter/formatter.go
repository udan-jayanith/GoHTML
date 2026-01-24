package Formatter

import (
	"bytes"
	"errors"
	"io"
	"strings"

	GoHtml "github.com/udan-jayanith/GoHTML"
	"golang.org/x/net/html"
)

// End of line
type EOL = string

const (
	// LF (Line Feed, \n): Standard for Unix, Linux.
	LF EOL = "\n"
	// CR (Carriage Return, \r): Legacy, used by pre-OS X Macs
	CR EOL = "\r"
	// CRLF (Carriage Return + Line Feed, \r\n): Standard for Windows
	CRLF EOL = "\r\n"
)

type FormatOptions struct {
	// Wether to wrap text
	// This does not effect text in attributes value.
	AutoWrapText bool
	// String used as indentation
	Tab string
	// Line ending
	LineEnding EOL
	// LineLength tries to maintain LineLength amount of bytes per line when formatting.
	LineLength int
}

var (
	defaultFormatOptions FormatOptions = FormatOptions{
		AutoWrapText: true,
		Tab:          "\t",
		LineEnding:   LF,
		LineLength:   70,
	}
)

type format_reader struct {
	tokenizer      *html.Tokenizer
	format_options FormatOptions
	// need a buf to hold processed strings for Reading.
	buf        bytes.Buffer
	last_token html.Token

	// indentation must be a empty string when initializing format_reader unless need margin to the left.
	indentation string
	// formatting buffer is a buffer use for string manipulation. This exists here to reduce memory reallocation.
	// use write_eol_to_formatting_buf if writing a eol line. Other wise use write_string_to_formatting_buf.
	// This is because format_reader needs to know is the last write is an eol.
	formatting_buf      bytes.Buffer
	is_last_write_eol   bool
	current_line_length int
}

func (fr *format_reader) increment_indentation() {
	fr.indentation += fr.format_options.Tab
}

func (fr *format_reader) decrement_indentation() {
	l := len(fr.indentation) - len(fr.format_options.Tab)
	if l < 0 {
		// TODO: Handle this with a error.
		return
	}
	fr.indentation = fr.indentation[:l]
}

// This interface is used for managing indentation and last eol.
type formatting_buf_writers interface {
	write_string_to_formatting_buf(str string)
	// write_eol_to_formatting_buf ends the line and add indentation.
	write_eol_to_formatting_buf()
	// try_write_new_line_to_formatting_buf does not tries to break str and if (line length) + len(str) is too long this will write str into a new line.
	try_write_new_line_to_formatting_buf(str string)
	// formatting_buf_line_length returns the current line length
	formatting_buf_line_length() int
	// formatting_buf_is_last_write_eol returns whether the last write used write_eol_to_formatting_buf.
	formatting_buf_is_last_write_eol() bool
}

// This also adds appropriate indentations.
func (fr *format_reader) write_eol_to_formatting_buf() {
	fr.formatting_buf.WriteString(fr.format_options.LineEnding)
	fr.formatting_buf.WriteString(fr.indentation)

	fr.is_last_write_eol = true
	fr.current_line_length = len(fr.format_options.LineEnding) + len(fr.indentation)
}

func (fr *format_reader) write_string_to_formatting_buf(str string) {
	fr.formatting_buf.WriteString(str)
	fr.is_last_write_eol = false
	fr.current_line_length += len(str)
}

func (fr *format_reader) try_write_new_line_to_formatting_buf(str string) {
	if fr.current_line_length+len(str) >= fr.format_options.LineLength {
		fr.write_eol_to_formatting_buf()
	}
	fr.write_string_to_formatting_buf(str)
}

func (fr *format_reader) formatting_buf_line_length() int {
	return fr.current_line_length
}

func (fr *format_reader) formatting_buf_is_last_write_eol() bool {
	return fr.is_last_write_eol
}

func (fr *format_reader) Read(b []byte) (n int, err error) {
	for {
		if fr.buf.Len() > 0 && len(b) > n {
			n1, _ := fr.buf.Read(b[n:])
			n += n1
		}

		if len(b) == n || n > 0 {
			return n, nil
		} else if fr.last_token.Type == html.ErrorToken && fr.buf.Len() == 0 {
			return n, io.EOF
		}

		switch fr.last_token.Type {
		// StartTagToken, SelfClosingTagToken, DoctypeToken and EndTagToken could have attributes
		// StartingTagToken increases indentation level
		// EndTagToken decreases indentation level
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken, html.DoctypeToken:
			tag_name := strings.ToLower(fr.last_token.Data)

			switch fr.last_token.Type {
			case html.EndTagToken:
				if !fr.formatting_buf_is_last_write_eol() {
					fr.write_eol_to_formatting_buf()
				}
				fr.write_string_to_formatting_buf("</" + tag_name)
			case html.DoctypeToken:
				tag_name = strings.ToUpper(tag_name)
				if !fr.formatting_buf_is_last_write_eol() {
					fr.write_eol_to_formatting_buf()
				}
				fr.write_string_to_formatting_buf("<" + tag_name)
			default:
				if !fr.formatting_buf_is_last_write_eol() {
					fr.write_eol_to_formatting_buf()
				}
				fr.write_string_to_formatting_buf("<" + tag_name)
			}

			format_attribute_list(tokenize_kv_attr(fr.last_token.Attr), fr.format_options.LineLength, fr)
			fr.write_string_to_formatting_buf(">")

			if fr.last_token.Type == html.StartTagToken {
				fr.increment_indentation()
			} else if fr.last_token.Type == html.EndTagToken {
				fr.decrement_indentation()
			}
		case html.TextToken:
			if !fr.format_options.AutoWrapText {
				break
			}
			// always write text in a new line and end in a new line.
			panic("Not implemented")
		case html.CommentToken:
			// always put comments on a new line and end in a new line.
		}

		fr.buf.Write(fr.formatting_buf.Bytes())
		fr.last_token = fr.tokenizer.Token()
	}
}

func (options *FormatOptions) Format(r io.Reader) (io.Reader, error) {
	if r == nil {
		return nil, errors.New("r is not optional")
	}

	t := html.NewTokenizer(r)
	return &format_reader{
		tokenizer:      t,
		format_options: *options,
		buf:            *bytes.NewBufferString(""),
		last_token:     t.Token(),
		formatting_buf: *bytes.NewBuffer(make([]byte, 0, 256)),
	}, nil
}

func (options *FormatOptions) FormatStrings(content string) (string, error) {
	return "", nil
}

func (options *FormatOptions) FormateFromString(content string) (io.Reader, error) {
	return nil, nil
}

func (options *FormatOptions) FormateFromNodeTree(node GoHtml.Node) (io.Reader, error) {
	return nil, nil
}

func (options *FormatOptions) FormatTo(r io.Reader, w io.Writer) error {
	return nil
}

func Format(r io.Reader) (io.Reader, error) {
	return defaultFormatOptions.Format(r)
}

func FormateFromNodeTree(node GoHtml.Node) (io.Reader, error) {
	return defaultFormatOptions.FormateFromNodeTree(node)
}

func FormatStrings(content string) (string, error) {
	return defaultFormatOptions.FormatStrings(content)
}
