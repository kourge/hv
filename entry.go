package main

import (
	"os"
	"io"
	"fmt"
	"crypto"
)

type Entry struct {
	Checksum string
	Filename string
}

func (e *Entry) String() string {
	return fmt.Sprintf("%s  %s", e.Checksum, e.Filename)
}

func (e *Entry) Calculate(h crypto.Hash) (sum []byte, err error) {
	if !h.Available() {
		return nil, HashUnavailableError{h}
	}
	hash := h.New()

	if file, err := os.Open(e.Filename); err != nil {
		return nil, err
	} else if _, err = io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func formatChecksum(h crypto.Hash, sum []byte) string {
	return fmt.Sprintf("%x", sum)
}

func (e *Entry) Fill(h crypto.Hash) (err error) {
	sum, err := e.Calculate(h)
	if err != nil {
		return
	}

	e.Checksum = formatChecksum(h, sum)
	return
}

func (e *Entry) Verify(h crypto.Hash) (match bool, err error) {
	sum, err := e.Calculate(h)
	if err != nil {
		return
	}

	match = formatChecksum(h, sum) == e.Checksum
	return
}

