package main

import (
	"os"
	"io"
	"fmt"
	"strings"
	"crypto"
	"container/list"
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

func matches(file os.FileInfo) bool {
	startsWith, endsWith := strings.HasPrefix, strings.HasSuffix
	name, mode := file.Name(), file.Mode()
	return mode.IsRegular() && !startsWith(name, ".") && !endsWith(name, "SUMS")
}

type Entries []*Entry

func EntriesFromFiles(files []os.FileInfo) Entries {
	items := 0
	entries := make([]*Entry, len(files))
	for _, file := range files {
		if matches(file) {
			entry := &Entry{Filename: file.Name()}
			entries[items] = entry
			items++
		}
	}

	return entries[0:items]
}

func EntriesFromChecksumFile(checksums *os.File) (entries Entries, err error) {
	entries = make([]*Entry, 0)
	r := NewReader(checksums)
	err = r.Each(func(entry *Entry) {
		entries = append(entries, entry)
	})
	return
}

type BucketsByChecksum map[string]*list.List

func (entries Entries) BucketsByChecksum() (buckets BucketsByChecksum) {
	buckets = make(map[string]*list.List)
	for _, entry := range entries {
		bucket, exists := buckets[entry.Checksum]
		if !exists {
			bucket = list.New()
			buckets[entry.Checksum] = bucket
		}
		bucket.PushBack(entry.Filename)
	}
	return
}

