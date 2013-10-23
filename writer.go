package main

import (
	"io"
	"fmt"
	"os"
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

func Dump(filename string, data map[string]string) (written int, err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}

	w := &Writer{file}
	for filename, checksum := range data {
		entry := &Entry{Filename: filename, Checksum: checksum}
		if n, err := w.WriteEntry(entry); err == nil {
			written += n
		} else {
			return written, err
		}
	}

	return
}

