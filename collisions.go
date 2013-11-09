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
	hashUsage := fmt.Sprintf("the hash to use; if unspecified, the following in tried in order: %s", strings.Join(preferredHashes, ", "))
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
		// Bucket files by size.
		switch groups := GroupBySize(bucket); {
		case len(groups) == bucket.Len():
			// All files are of different lengths. Then none of them are identical.
			fileinfos := list.New()
			for _, group := range groups {
				fileinfos.PushBack(group.Front().Value.(os.FileInfo))
			}
			emit(checksum, fileinfos)
		case len(groups) == 1:
			// All files are of the same length. It is possible that they are all
			// identical, but further checks are needed.
		default:
			// Some files are of the same length. Verify each set of them that
			// share lengths and check if they are identical.
		}
	}
}

func emit(checksum string, fileinfos *list.List) {
	warn("%s\n", checksum)
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
		group.PushBack(info)
	}
	return
}

