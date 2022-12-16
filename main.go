package main

import (
	"os"

	"github.com/braswelljr/rmx/commands"
)

func main() {
	command := commands.RootCommand()

	// add flags
	commands.CommandFlags(command, &commands.Rmc)

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
