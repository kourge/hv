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
	Run: collisions,
	Usage: `collisions [-c=hash] [-D=dir]`,
	Short: "Find hash collisions within a checksum file",
	Long: `
Find all instances of hash collisions for the given directory's checksum file, i.e.
display all sets of files that have the same checksum but are not content-identical.`,
}

// Initialized in common:
// var preferredHashes []string

func init() {
	const (
		cwdUsage = "the directory for which to verify using the checksum file"
	)
	hashUsage := fmt.Sprintf("the hash to use; if unspecified, the following are tried in order: %s", strings.Join(preferredHashes, ", "))
	f := &cmdCollisions.Flag
	f.Var(&hashFunction, "c", hashUsage)
	f.StringVar(&cwd, "D", ".", cwdUsage)
}

func collisions(cmd *Command, args []string) {
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
				warn("%s\n", err)
			}
		}

		fileinfos := list.New()
		for _, group := range groups {
			fileinfos.PushBack(group.Front().Value.(os.FileInfo))
		}

		switch {
		case len(groups) == bucket.Len():
			// All files are of different lengths. Then none of them are identical,
			// and all files under this checksum are genuinely colliding. This is
			// a plain collision.
			emitPlain(checksum, fileinfos)
		default:
			// Some or all of the files are of the same length. Some of them might
			// be genuinely identical.
			warnSameLen(checksum, fileinfos)
		}
	}
}

func emitPlain(checksum string, fileinfos *list.List) {
	warn("%s\n", checksum)
	for e := fileinfos.Front(); e != nil; e = e.Next() {
		fileinfo := e.Value.(os.FileInfo)
		warn("\t%s\n", fileinfo.Name())
	}
	warn("\n")
}

func warnSameLen(checksum string, fileinfos *list.List) {
	warn("%s\n", checksum)
	warn("Some or all of the following files are of the same size but hash to the same checksum above\n")
	for e := fileinfos.Front(); e != nil; e = e.Next() {
		fileinfo := e.Value.(os.FileInfo)
		warn("\t%s\n", fileinfo.Name())
	}
	warn("\n")
}

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

