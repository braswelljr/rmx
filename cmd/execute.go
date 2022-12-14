package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/cmd/help"
	"github.com/braswelljr/rmx/cmd/interactive"
	"github.com/braswelljr/rmx/rm"
)

// Run - runs the command.
//
//	@param command - commands to run.
//	@param args - arguments and flags to run alongside the command.
//	@return error - error if there is one.
func Run(r *rm.Rm, command *cobra.Command, args []string) error {
	// check for arguments with no flags
	// if len(args) > 0 && command.Flags().NFlag() == 0 {
	// 	// execute the command

	// }

	// check for commands and execute them accordingly
	// check for help command
	// check for empty args or help flag
	if command.Flags().Changed("help") || len(args) < 1 {
		// run the help command
		help.Help(r, command, args)
		return nil
	}

	// check for interactive flag
	if r.Ii || command.Flags().Changed("interactive") {
		// run the interactive command
		interactive.InteractiveIi(args)
	}

	// return on no error
	return nil
}

// ArgumentValidator - Validates the given arguments
func ArgumentValidator(rmx *rm.Rm) func(command *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		cmd_name := color.CyanString("rmx")
		// check if args are empty
		if len(args) < 1 {
			// print prompt
			fmt.Printf("%s: missing arguments or flags for command\nTry '%s --help' for more information.\n", cmd_name, cmd_name)
			return nil
		}

		// return on no error
		return nil
	}
}
