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

func (r *Reader) Iter() chan *Entry {
	ch := make(chan *Entry)
	go func() {
		for {
			if entry, err := r.ReadEntry(); err != nil {
				break
			} else if entry != nil {
				ch <- entry
			}
		}
		close(ch)
	}()

	return ch
}

func Load(filename string) (m map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	r := NewReader(file)
	m = make(map[string]string)
	for {
		if e, err := r.ReadEntry(); err == io.EOF {
			break
		} else if err != nil {
			return m, err
		} else if e != nil {
			m[e.Filename] = e.Checksum
		}
	}

	return m, nil
}

