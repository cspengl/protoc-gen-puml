package generator

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type IndentWriter struct {
	Indent  string
	Wrapped *protogen.GeneratedFile
	level   int
}

func (w *IndentWriter) Write(p []byte) (n int, err error) {
	n_indent, err := w.Wrapped.Write(w.indent())
	if err != nil {
		return n_indent, err
	}
	n_content, err := w.Wrapped.Write(p)
	return n_indent + n_content, err
}

func (w *IndentWriter) P(content ...any) {
	w.Wrapped.P(append([]any{string(w.indent())}, content...)...)
}

func (w *IndentWriter) indent() []byte {
	return []byte(strings.Repeat(w.Indent, w.level))
}

func (w *IndentWriter) IncrLevel() {
	w.level += 1
}

func (w *IndentWriter) DecrLevel() {
	w.level -= 1
}
