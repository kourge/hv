package main

import (
	"os"
)

var cmdHelp = &Command{
	Usage: `help <command>`,
	Short: "Show usage for a command",
	Long: `
Display a detailed view of a command's usage.`,
}

func init() {
	// avoid initialization loop from specifying this as part of the struct literal
	cmdHelp.Run = runHelp
}

func PrintCommands() {
	warn(`Available commands:

`)
	for _, command := range commands {
		warn("  %s - %s\n", command.Name(), command.Short)
	}
}

func helpUsage(cmd *Command) {
	cmd.PrintUsage(nil)
	PrintCommands()
	os.Exit(1)
}

func runHelp(cmd *Command, args []string) {
	if len(args) == 0 {
		helpUsage(cmd)
	}

	var command *Command
	topic := args[0]
	for _, cmd := range commands {
		if cmd.Name() == topic {
			command = cmd
			break
		}
	}

	if command != nil {
		command.PrintUsage(nil)
	} else {
		warn("Could not find the command \"%s\".\n", topic)
		helpUsage(cmd)
	}
}

