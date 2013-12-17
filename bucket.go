package main

import (
	"bytes"
)

type Bucket struct {
	Files []*LazyFile
}

func (b *Bucket) String() string {
	var output bytes.Buffer
	output.WriteString("{")
	size := len(b.Files)
	for i, file := range b.Files {
		output.WriteString(file.String())
		if i != size - 1 {
			output.WriteString(", ")
		}
	}
	output.WriteString("}")
	return output.String()
}

func (b *Bucket) Contains(file *LazyFile) (contains bool, err error) {
	if b.Files == nil {
		return false, nil
	}

	// Every file in a Bucket is identical, so only one comparison is needed.
	return b.Files[0].Equal(file)
}

func (b *Bucket) Add(file *LazyFile) {
	b.Files = append(b.Files, file)
}

type Buckets []*Bucket

type LazyFileBucketer struct {
	Buckets
}

func (b *LazyFileBucketer) Add(file *LazyFile) error {
	for _, bucket := range b.Buckets {
		contains, err := bucket.Contains(file)
		if err != nil {
			return err
		} else if contains {
			bucket.Add(file)
			return nil
		}
	}

	// The file being added is not identical to any known bucket.
	bucket := &Bucket{}
	bucket.Add(file)
	b.Buckets = append(b.Buckets, bucket)
	return nil
}

