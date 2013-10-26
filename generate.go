package main

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
)

var (
	force bool
	hashFunction HashValue
	cwd string
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
	hash := &hashFunction
	if hash.String() == "" {
		hash.Set("SHA1")
	}

	if dirs, err := ioutil.ReadDir(cwd); err == nil {
		for _, file := range dirs {
			if matches(file) {
				entry := &Entry{Filename: file.Name()}
				if err := entry.Fill(hash.Hash); err == nil {
					warn("%s\n", entry)
				} else {
					croak(err)
				}
			}
		}
	} else {
		croak(err)
	}
}

