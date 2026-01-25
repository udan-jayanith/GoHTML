package Formatter

import (
	"bytes"
	"errors"
	"io"

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
