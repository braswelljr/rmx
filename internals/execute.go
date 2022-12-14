package internals

import "github.com/spf13/cobra"

// Execute - is the entry point for the Cobra CLI
// @param command - the root command
// @param args - the arguments passed to the CLI
// @return error - the error if any
func Execute(command *cobra.Command, args []string) error {
	// set the command
	command.SetArgs(args)
	// execute the command
	return command.Execute()
}
