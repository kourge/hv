package main

import (
	"os"
	"errors"
)

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

func setDirAndHashOptions() (err error, hash *HashValue, checksums *os.File) {
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

