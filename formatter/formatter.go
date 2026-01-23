package Formatter

import (
	"io"

	GoHtml "github.com/udan-jayanith/GoHTML"
)

// End of line
type EOL = string

const (
	// CRLF (Carriage Return + Line Feed, \r\n): Standard for Windows
	CRLF EOL = "\r\n"
	// LF (Line Feed, \n): Standard for Unix, Linux.
	LF EOL = "\n"
	// CR (Carriage Return, \r): Legacy, used by pre-OS X Macs
	CR EOL = "\r"
)

type FormatOptions struct {
	AutoWrap struct {
		Text, HTML, CommentText bool
	}
	Tab       string
	LineEnding EOL
}

var (
	DefaultFormatOptions FormatOptions = FormatOptions{
		AutoWrap: struct {
			Text        bool
			HTML        bool
			CommentText bool
		}{
			Text: true,
			HTML: true,
		},
	}
)

func (options *FormatOptions) Format(r io.Reader) io.Reader {
	return nil
}

func (options *FormatOptions) FormatStrings(content string) string {
	return ""
}

func (options *FormatOptions) FormateFromString(content string) io.Reader {
	return nil
}

func (options *FormatOptions) FormateFromNodeTree(node GoHtml.Node) io.Reader {
	return nil
}

func Format(r io.Reader) io.Reader {
	return DefaultFormatOptions.Format(r)
}

func FormateFromNodeTree(node GoHtml.Node) io.Reader {
	return DefaultFormatOptions.FormateFromNodeTree(node)
}

func FormatStrings(content string) string {
	return DefaultFormatOptions.FormatStrings(content)
}
