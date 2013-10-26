package main

import (
	"os"
	"flag"
	"fmt"
)

func warn(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func croak(e error) (n int, err error) {
	return fmt.Fprintf(os.Stderr, "%s\n", e)
}

func usage() {
	program := os.Args[0]
	warn("usage: %s [SUMS]\n", program)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	filename := args[0]

	if result, err := Load(filename); err == nil {
		fmt.Printf("result = %#v\n", result)
	} else {
		warn("error: %v\n", err)
	}
}
