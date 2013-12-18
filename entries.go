package main

import (
	"os"
	"io/ioutil"
	"container/list"
	"strings"
)

type Entries []*Entry

func EntriesFromPath(path string) (entries Entries, err error) {
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		return
	}
	return EntriesFromFiles(files), nil
}

func matches(file os.FileInfo) bool {
	startsWith, endsWith := strings.HasPrefix, strings.HasSuffix
	name, mode := file.Name(), file.Mode()
	return mode.IsRegular() && !startsWith(name, ".") && !endsWith(name, "SUMS")
}

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

