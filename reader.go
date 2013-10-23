package main

import (
	"os"
	"fmt"
	"io"
	"strings"
	"bufio"
)
var _ = fmt.Println

const Separator = "  "

type Reader struct {
	*bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

func (r *Reader) ReadEntry() (entry *Entry, err error) {
	bytes, _, err := r.ReadLine()
	if err != nil {
		return
	}

	line := string(bytes)
	sep := strings.Index(line, Separator)
	if sep == -1 {
		return
	}

	entry = &Entry{Checksum: line[0:sep], Filename: line[sep+len(Separator):]}
	return
}

func (r *Reader) Each(f func(*Entry)) (err error) {
	for {
		if e, err := r.ReadEntry(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if e != nil {
			f(e)
		}
	}
	return
}

func Load(filename string) (m map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	m = make(map[string]string)
	r := NewReader(file)
	err = r.Each(func(e *Entry) {
		m[e.Filename] = e.Checksum
	})

	return m, err
}

