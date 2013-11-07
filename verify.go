package main

import (
	"fmt"
	"strings"
	"os"
)

var (
	silent bool
	// Declared in common:
	// hashFunction HashValue
	// cwd string
)

var cmdVerify = &Command{
	Run: verify,
	Usage: `verify [-s] [-c=hash] [-D=dir]`,
	Short: "Verify using a checksum file",
	Long: `
Verify all files for the given directory against a checksum file.`,
}

// Initialized in common:
// var preferredHashes []string

func init() {
	const (
		silentUsage = "silent; don't output to STDERR"
		cwdUsage = "the directory for which to verify using the checksum file"
	)
	hashUsage := fmt.Sprintf("the hash to use; if unspecified, the following in tried in order: %s", strings.Join(preferredHashes, ", "))
	f := &cmdVerify.Flag
	f.BoolVar(&silent, "s", false, silentUsage)
	f.Var(&hashFunction, "c", hashUsage)
	f.StringVar(&cwd, "D", ".", cwdUsage)
}

func verify(cmd *Command, args []string) {
	err, hash, checksums := setDirAndHashOptions()

	allMatch := true
	r := NewReader(checksums)
	err = r.Each(func(entry *Entry) {
		if ok, err := entry.Verify(hash.Hash); err != nil {
			if !silent {
				warn("%s\n", err)
			}
			allMatch = false
		} else if !ok {
			if !silent {
				warn("%s does not match %s\n", entry.Filename, entry.Checksum)
			}
			allMatch = false
		}
	})
	if err != nil {
		die(err)
	}

	if allMatch {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

