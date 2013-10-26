package main

import (
	"os"
	"io"
	"fmt"
	"flag"
	"strings"
)

type Command struct {
	Run func(cmd *Command, args []string)
	Flag flag.FlagSet

	Usage string
	Short string
	Long string
}

func (c *Command) Name() (name string) {
	name = c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) PrintUsage(w io.Writer) {
	program := os.Args[0]

	if w == nil {
		w = os.Stderr
	}
	if c.Runnable() {
		fmt.Fprintf(w, "Usage: %s %s\n", program, c.Usage)
		fmt.Fprintf(w, "%s\n", c.Short)
		if c.Flag.NFlag() != 0 {
			fmt.Fprint(w, "\n")
		}
		c.Flag.PrintDefaults()
	}
	fmt.Fprintln(w, strings.TrimRight(c.Long, "\n"))
}

