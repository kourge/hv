package main

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
)

var (
	force bool
	// Declared in main:
	// hashFunction HashValue
	// cwd string
)

var cmdGenerate = &Command{
	Run: generate,
	Usage: `generate [-f] [-c=hash] [-D=dir]`,
	Short: "Generate a checksum file",
	Long: `
Generate a checksum file for the given directory. The generated checksum file
will be named %sSUMS, where %s is the chosen hash function name in all caps.`,
}

func init() {
	const (
		forceUsage = "overwrite an existing checksum file"
		cwdUsage = "the directory for which to generate the checksum file"
	)
	hashUsage := fmt.Sprintf("the hash function to use, e.g. %s", strings.Join(hashFunction.Values(), ", "))
	f := &cmdGenerate.Flag
	f.BoolVar(&force, "f", false, forceUsage)
	f.Var(&hashFunction, "c", hashUsage)
	f.StringVar(&cwd, "D", ".", cwdUsage)
}

func generate(cmd *Command, args []string) {
	if err := os.Chdir(cwd); err != nil {
		die(err)
	}

	hash := &hashFunction
	if hash.String() == "" {
		hash.Set("SHA1")
	}

	checksumFile := hash.Filename()
	if _, err := os.Stat(checksumFile); err == nil && !force {
		warn("%s already exists\n", checksumFile)
		os.Exit(1)
	}
	f, err := os.Create(checksumFile)
	if err != nil {
		die(err)
	}

	w := NewWriter(f)

	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		die(err)
	}

	entries := EntriesFromFiles(files)
	for _, entry := range entries {
		if err := entry.Fill(hash.Hash); err != nil {
			die(err)
		}
	}

	for _, entry := range entries {
		w.WriteEntry(entry)
	}

	os.Exit(0)
}

