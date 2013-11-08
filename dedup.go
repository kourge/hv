package main

import (
	"fmt"
	"strings"
	"os"
	"container/list"
)

var (
	// Declared in common:
	// hashFunction HashValue
	// cwd string
	dryRun bool
)

var cmdDedup = &Command{
	Run: dedup,
	Usage: `dedup [-c=hash] [-D=dir] [--dryrun]`,
	Short: "deduplicate using a checksum file",
	Long: `
Deduplicate all files for the given directory against a checksum file.`,
}

// Initialized in common:
// var preferredHashes []string

func init() {
	const (
		cwdUsage = "the directory for which to verify using the checksum file"
		dryRunUsage = "only output what would have been done; do not perform any destructive operations"
	)
	hashUsage := fmt.Sprintf("the hash to use; if unspecified, the following in tried in order: %s", strings.Join(preferredHashes, ", "))
	f := &cmdDedup.Flag
	f.Var(&hashFunction, "c", hashUsage)
	f.StringVar(&cwd, "D", ".", cwdUsage)
	f.BoolVar(&dryRun, "dryrun", false, dryRunUsage)
}

func dedup(cmd *Command, args []string) {
	_, checksums := setDirAndHashOptions()

	entries, err := EntriesFromChecksumFile(checksums)
	if err != nil {
		die(err)
	}

	if dryRun {
		warn("# Dry run mode is on\n")
	}

	buckets := entries.BucketsByChecksum()
	for checksum, bucket := range buckets {
		if bucket.Len() > 1 {
			promptForRemoval(bucket, checksum)
		}
	}
}

func promptForRemoval(duplicates *list.List, checksum string) {
PROMPT:
	for e, i := duplicates.Front(), 1; e != nil; e, i = e.Next(), i+1 {
		file := e.Value.(string)
		warn("# [%d] %s\n", i, file)
	}
	warn("# All of these have checksum %s. Keep which? ", checksum)

	var choice int
	if n, err := fmt.Scanf("%d", &choice); err == nil && n == 1 {
		if choice < 1 || choice > duplicates.Len() {
			warn("# %d is not a valid choice\n\n", choice)
			goto PROMPT
		}
		removeDuplicatesAndKeep(duplicates, choice)
	} else {
		warn("# Please enter a number\n\n")
		goto PROMPT
	}
}

func removeDuplicatesAndKeep(duplicates *list.List, choice int) {
	e := duplicates.Front()
	for i := 1; i < choice; i++ {
		e = e.Next()
	}

	file := duplicates.Remove(e).(string)
	warn("# Keeping %s\n", file)

	for e = duplicates.Front(); e != nil; e = e.Next() {
		file := e.Value.(string)

		if !dryRun {
			if err := os.Remove(file); err != nil {
				croak(err)
			}
		}
		warn("rm %s\n", file)
	}
	warn("\n")
}

