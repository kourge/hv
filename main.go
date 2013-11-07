package main

import (
	"os"
	"strings"
)

var commands = []*Command {
	cmdGenerate,
	cmdVerify,
	cmdDedup,
	cmdHelp,
}

func usage() {
	program := os.Args[0]
	message := strings.TrimLeft(`
usage: %s <command>

`, "\n")
	warn(message, program)

	PrintCommands()
	os.Exit(1)
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
