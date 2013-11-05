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

func matches(file os.FileInfo) bool {
	startsWith, endsWith := strings.HasPrefix, strings.HasSuffix
	name, mode := file.Name(), file.Mode()
	return mode.IsRegular() && !startsWith(name, ".") && !endsWith(name, "SUMS")
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

	dirs, err := ioutil.ReadDir(cwd)
	if err != nil {
		die(err)
	}

	for _, file := range dirs {
		if matches(file) {
			entry := &Entry{Filename: file.Name()}
			if err := entry.Fill(hash.Hash); err == nil {
				w.WriteEntry(entry)
			} else {
				die(err)
			}
		}
	}

	os.Exit(0)
}

