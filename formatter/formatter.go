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
	buf            bytes.Buffer
	last_token     html.Token
	
	indentation string
	// formatting buffer is a buffer use for string manipulation. This exists here to reduce memory reallocation.
	formatting_buf bytes.Buffer
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
		// StartTagToken, SelfClosingTagToken, DoctypeToken and EndTagToken cloud have attributes
		// StartingTagToken increases indentation level
		// EndTagToken decreases indentation level
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken, html.DoctypeToken:
			tag_name := strings.ToLower(fr.last_token.Data)
			format_attribute_list(tokenize_kv_attr(fr.last_token.Attr), &fr.format_options, "\t", &fr.formatting_buf)
			
			switch fr.last_token.Type {
			case html.StartTagToken:
			case html.SelfClosingTagToken:
			case html.EndTagToken:
			case html.DoctypeToken:
				// Make doctype tags tag name uppercase
			}
		case html.TextToken:
			// Wrap the text if specified
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
		indentation: options.Tab,
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
