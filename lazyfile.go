package main

import (
	"os"
	"bytes"
	"fmt"
)

const CHUNK_SIZE int64 = 4096

type LazyFile struct {
	os.FileInfo
	file *os.File
	data []byte
}

func (fi *LazyFile) String() string {
	return fmt.Sprintf("{\"%s\", %d}", fi.FileInfo.Name(), len(fi.data))
}

func (fi *LazyFile) Open() (err error) {
	if fi.file != nil {
		return
	}
	file, err := os.Open(fi.Name())
	if err == nil {
		fi.file = file
	}
	return
}

func (fi *LazyFile) Close() (err error) {
	if fi.file == nil {
		return nil
	}
	return fi.file.Close()
}

func (fi *LazyFile) Slurp() (n int, err error) {
	if err := fi.Open(); err != nil {
		return 0, err
	}
	buffer := make([]byte, CHUNK_SIZE)
	n, err = fi.file.Read(buffer)
	if err == nil {
		fi.data = append(fi.data, buffer...)
	}
	return
}

func (fi *LazyFile) Equal(other *LazyFile) (equal bool, err error) {
	// Check both files' length first
	if fi.Size() != other.Size() {
		return false, nil
	}
	size := fi.Size()

	var i int64 = 0
	for i < size {
		// Slurp when necessary
		if i + CHUNK_SIZE > int64(len(fi.data)) {
			if _, err := fi.Slurp(); err != nil {
				return false, err
			}
		}
		if i + CHUNK_SIZE > int64(len(other.data)) {
			if _, err := other.Slurp(); err != nil {
				return false, err
			}
		}

		x, y := fi.data[i:i+CHUNK_SIZE], other.data[i:i+CHUNK_SIZE]
		if !bytes.Equal(x, y) {
			return false, nil
		}
		i += CHUNK_SIZE
	}

	return true, nil
}

