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
	// CharLenPerLine tries to maintain CharLenPerLine amount of bytes per line when formatting.
	CharLenPerLine int
}

var (
	defaultFormatOptions FormatOptions = FormatOptions{
		AutoWrapText: true,
		Tab:          "\t",
		LineEnding:   LF,
	}
)

type format_reader struct {
	tokenizer      *html.Tokenizer
	format_options FormatOptions
	// need a buf to hold processed strings for Reading.
	buf            bytes.Buffer
	last_token     html.Token

	// indentation must be a empty string when initializing format_reader unless need margin to the left.  
	indentation string
	// formatting buffer is a buffer use for string manipulation. This exists here to reduce memory reallocation.
	// use write_eol_to_formatting_buf if writing a eol line. Other wise use write_string_to_formatting_buf.
	// This is because format_reader needs to know is the last write is an eol. 
	formatting_buf bytes.Buffer
	is_last_write_eol bool
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

func (fr *format_reader) write_eol_to_formatting_buf() {
	fr.formatting_buf.WriteString(fr.format_options.LineEnding)
	fr.is_last_write_eol = true
}

func (fr *format_reader) write_string_to_formatting_buf(str string) {
	fr.formatting_buf.WriteString(str)
	fr.is_last_write_eol = false
}

func (fr *format_reader) Read(b []byte) (n int, err error) {
	for {
		if fr.buf.Len() > 0 {
			n1, _ := fr.buf.Read(b[n:])
			n += n1
		}

		if len(b) == n || fr.buf.Len() > 0 && n > 0 {
			return n, nil
		} else if fr.last_token.Type == html.ErrorToken && fr.buf.Len() == 0 {
			return n, io.EOF
		}

		var str string
		switch fr.last_token.Type {
		// StartTagToken, SelfClosingTagToken, DoctypeToken and EndTagToken could have attributes
		// StartingTagToken increases indentation level
		// EndTagToken decreases indentation level
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken, html.DoctypeToken:
			tag_name := strings.ToLower(fr.last_token.Data)

			switch fr.last_token.Type {
			case html.EndTagToken:
				fr.formatting_buf.WriteString("</" + tag_name)
				fr.decrement_indentation()
			case html.DoctypeToken:
				tag_name = strings.ToUpper(tag_name)
				fr.formatting_buf.WriteString("<" + tag_name)
			default:
				fr.formatting_buf.WriteString("<" + tag_name)
			}

			if len(fr.last_token.Attr) > 0 {
				fr.formatting_buf.WriteString(" ")
				format_attribute_list(tokenize_kv_attr(fr.last_token.Attr), &fr.format_options, fr.indentation, &fr.formatting_buf)
			}

			if fr.last_token.Type == html.StartTagToken {
				fr.increment_indentation()
			}
			fr.formatting_buf.WriteString(">")
		case html.TextToken:
			if !fr.format_options.AutoWrapText {
				break
			}
			panic("Not implemented")
		case html.CommentToken:
			//Alway put comments on a new line
		}

		_, err := fr.buf.WriteString(str)
		if err != nil {
			return 0, err
		}
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
