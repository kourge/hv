package main

import (
	"fmt"
	"os"
	"errors"
)

func warn(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func croak(e error) {
	warn("%s\n", e)
}

func die(e error) {
	croak(e)
	os.Exit(1)
}

var (
	hashFunction HashValue
	cwd string
	dryRun bool
	silent bool
)

var preferredHashes = []string{"SHA512", "SHA1", "MD5"}

func findChecksumFile() (hash *HashValue, file *os.File, err error) {
	for _, tryHash := range preferredHashes {
		hash := &HashValue{}
		hash.Set(tryHash)
		file, err := os.Open(hash.Filename())
		if err == nil {
			return hash, file, err
		}
	}
	return nil, nil, errors.New("No known checksum files found in directory")
}

func setDirAndHashOptions() (hash *HashValue, checksums *os.File) {
	var err error

	if err := os.Chdir(cwd); err != nil {
		die(err)
	}

	if hashFunction.Hash == 0x0 {
		hash, checksums, err = findChecksumFile()
		if err != nil {
			die(err)
		}
	} else {
		hash = &hashFunction
		checksums, err = os.Open(hash.Filename())
		if err != nil {
			die(err)
		}
	}

	return
}

