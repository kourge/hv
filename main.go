package main

import (
	"os"
	"fmt"
	"strings"
)

func warn(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func croak(e error) (n int, err error) {
	return fmt.Fprintf(os.Stderr, "%s\n", e)
}

func die(e error) {
	croak(e)
	os.Exit(1)
}

var commands = []*Command {
	cmdGenerate,
}

func usage() {
	program := os.Args[0]
	message := strings.TrimLeft(`
usage: %s <command>

Available commands:

`, "\n")
	warn(message, program)

	for _, command := range commands {
		warn("  %s - %s\n", command.Name(), command.Short)
	}
	os.Exit(2)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}

	args := os.Args[1:]
	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.PrintUsage(nil)
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	warn("Unknown command: %s\n", args[0])
	usage()
}
