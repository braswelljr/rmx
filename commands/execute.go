package commands

import (
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/commands/help"
	"github.com/braswelljr/rmx/rm"
)

// Run - runs the command.
//
//	@param command - commands to run.
//	@param args - arguments and flags to run alongside the command.
//	@return error - error if there is one.
func Run(r *rm.RM, command *cobra.Command, args []string) error {
	// check for empty args or help flag
	if command.Flags().NFlag() < 1 || command.Flags().Changed("help") {
		return help.Help(r, command, args)
	}

	// check for commands and execute them accordingly

	// return on no error
	return nil
}

// ArgumentValidator - Validates the given arguments
func ArgumentValidator(rmc *rm.RM) func(command *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		// check if args are empty
		if len(args) < 1 {
			println("rmx: missing arguments or flags for command")
			println("Try 'rmx --help' for more information.")
			return nil
		}

		// return on no error
		return nil
	}
}
