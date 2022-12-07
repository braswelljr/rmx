package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// initialize cobra
	command := &cobra.Command{
		Use:   "rmx [OPTION]... [FILE]...",
		Short: "A cross platform drop-in replacement for rm",
		Run: func(command *cobra.Command, args []string) {
			// check for help flag
		},
	}

	// add flags
	command.PersistentFlags().BoolP("help", "h", false, "help for this command")

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
