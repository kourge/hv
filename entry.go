package main

import (
	"os"
	"io"
	"fmt"
	"crypto"
)

type HashUnavailableError struct {
	h crypto.Hash
}

func (e *HashUnavailableError) Error() string {
	var hashName string
	switch e.h {
	case crypto.MD4: hashName = "MD4"
	case crypto.MD5: hashName = "MD5"
	case crypto.SHA1: hashName = "SHA1"
	case crypto.SHA224: hashName = "SHA224"
	case crypto.SHA256: hashName = "SHA256"
	case crypto.SHA384: hashName = "SHA384"
	case crypto.SHA512: hashName = "SHA512"
	case crypto.MD5SHA1: hashName = "MD5SHA1"
	case crypto.RIPEMD160: hashName = "RIPEMD160"
	default: hashName = "unknown hash"
	}
	return fmt.Sprint("%s is not a known hash function or not linked into the binary", hashName)
}

type Entry struct {
	Checksum string
	Filename string
}

func (e *Entry) String() string {
	return fmt.Sprintf("%s  %s", e.Checksum, e.Filename)
}

func (e *Entry) Calculate(h crypto.Hash) (sum []byte, err error) {
	if !h.Available() {
		return nil, &HashUnavailableError{h}
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

