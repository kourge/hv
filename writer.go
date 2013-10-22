package main

import (
	"io"
	"fmt"
)

type Writer struct {
	io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

func (w *Writer) WriteEntry(entry *Entry) (written int, err error) {
	return fmt.Fprintf(w, "%s\n", entry.String())
}

func (e *Entry) Write(dst io.Writer) (written int, err error) {
	w := &Writer{dst}
	return w.WriteEntry(e)
}

