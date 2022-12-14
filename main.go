package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// initialize cobra
	command := &cobra.Command{
		Use:   "rmx [OPTION]... [FILE]...",
		Short: "A cross-platform replacement for UNIX rm command",
		Long:  "A cross-platform replacement for UNIX rm command",
		Run: func(command *cobra.Command, args []string) {
			// Execute the command
			command.SetArgs(args)
		},
	}

	// add flags
	command.PersistentFlags().BoolP("help", "h", false, "help for this command")

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
