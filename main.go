package main

import (
	"fmt"
	"os"

	commands "github.com/braswelljr/rmx/cmd"
)

func main() {
	// get root command
	command := commands.RootCommand()

	// add flags
	commands.CommandFlags(command, &commands.Rmc)

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err) // print error to stderr
		os.Exit(1)                   // exit on error
	}
}
