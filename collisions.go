package main

import (
	"fmt"
	"strings"
	"container/list"
	"os"
)

var (
	// Declared in common:
	// hashFunction HashValue
	// cwd string
)

var cmdCollisions = &Command{
	Run: runCollisions,
	Usage: `collisions [-c=hash] [-D=dir]`,
	Short: "Find hash collisions within a checksum file",
	Long: `
Find all instances of hash collisions for the given directory's checksum file, i.e.
display all sets of files that have the same checksum but are not content-identical.
Under a particular checksum, files that are actually identical will be clustered
together, while files that are distinct are separated by a blank line.

Expensive operations are avoided: if two files share the same checksum but are of
different length, they are considered to be colliding and their contents are not
actually compared.

Warning: if your checksum file is not up-to-date or is incorrect, then the result
of running this will be incorrect.`,
}

// Initialized in common:
// var preferredHashes []string

func init() {
	const (
		cwdUsage = "the directory in which to look for collisions using its checksum file"
	)
	hashUsage := fmt.Sprintf("the hash to use; if unspecified, the following are tried in order: %s", strings.Join(preferredHashes, ", "))
	f := &cmdCollisions.Flag
	f.Var(&hashFunction, "c", hashUsage)
	f.StringVar(&cwd, "D", ".", cwdUsage)
}

func runCollisions(cmd *Command, args []string) {
	_, checksums := setDirAndHashOptions()

	entries, err := EntriesFromChecksumFile(checksums)
	if err != nil {
		die(err)
	}

	buckets := entries.BucketsByChecksum()
	for checksum, bucket := range buckets {
		// A single file for a given checksum indicates no collisions
		if bucket.Len() == 1 {
			continue
		}

		// Bucket files by size.
		groups := GroupBySize(bucket)
		if group, exists := groups[-1]; exists {
			for e := group.Front(); e != nil; e = e.Next() {
				err := e.Value.(error)
				croak(err)
			}
		}

		switch {
		case len(groups) == bucket.Len():
			// All files are of different lengths. Then none of them are identical,
			// and all files under this checksum are genuinely colliding. This is
			// a plain collision.
			emitPlain(checksum, groups)
		default:
			// Some or all of the files are of the same length. Some of them might
			// be genuinely identical.
			emitMixed(checksum, groups)
		}
	}
}

func emitPlain(checksum string, groups map[int64]*list.List) {
	warn("%s\n", checksum)
	for _, group := range groups {
		fileinfo := group.Front().Value.(os.FileInfo)
		warn("\t%s\n\n", fileinfo.Name())
	}
}

func emitMixed(checksum string, groups map[int64]*list.List) {
	warn("%s\n", checksum)
	for _, group := range groups {
		if group.Len() == 1 {
			fileinfo := group.Front().Value.(os.FileInfo)
			warn("\t%s\n\n", fileinfo.Name())
			continue
		}

		buckets, errs := GroupByContent(group)
		for _, bucket := range buckets {
			for _, lazyfile := range bucket.Files {
				warn("\t%s\n", lazyfile.Name())
			}
			warn("\n")
		}
		for _, err := range errs {
			warn("\t%s\n", err)
		}
	}
}

// Group a list of strings representing files by their stat()ed file sizes. The
// return value is a map, where the key is a file size and the value is a
// list.List of os.FileInfo. For any original file that could not be stat()ed,
// the returned error is added to the list.List corresponding to the invalid
// size of -1.
func GroupBySize(entries *list.List) (groups map[int64]*list.List) {
	groups = make(map[int64]*list.List)
	for e := entries.Front(); e != nil; e = e.Next() {
		var size int64 = -1
		filename := e.Value.(string)

		info, err := os.Lstat(filename)
		if err == nil {
			size = info.Size()
		}

		group, exists := groups[size]
		if !exists {
			group = list.New()
			groups[size] = group
		}

		if size != -1 {
			group.PushBack(info)
		} else {
			group.PushBack(err)
		}
	}
	return
}

// Segment a list of strings representing files by their content treated as a
// byte array. Returns an array of pointers to a Bucket, each of which contain
// an array of pointers to a LazyFile. All LazyFiles within the same bucket are
// byte-by-byte identical.
func GroupByContent(entries *list.List) (buckets Buckets, errs []error) {
	if entries.Len() == 0 {
		return nil, nil
	}

	b := &LazyFileBucketer{}
	for e := entries.Front(); e != nil; e = e.Next() {
		fileinfo := e.Value.(os.FileInfo)
		file := &LazyFile{FileInfo: fileinfo}
		if err := b.Add(file); err != nil {
			errs = append(errs, err)
		}
	}

	return b.Buckets, errs
}

