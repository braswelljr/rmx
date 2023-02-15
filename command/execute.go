package command

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/command/directory"
	"github.com/braswelljr/rmx/command/help"
	"github.com/braswelljr/rmx/command/interactive"
	"github.com/braswelljr/rmx/command/verbose"
	"github.com/braswelljr/rmx/rm"
)

// Run - runs the command.
//
//	@param command - commands to run.
//	@param args - arguments and flags to run alongside the command.
//	@return error - error if there is one.
func Run(r *rm.Rm, command *cobra.Command, args []string) error {
	// check for arguments with no flags
	if len(args) > 0 && command.Flags().NFlag() == 0 {
		// execute the command
		if err := rm.RemoveMultiple(args); err != nil {
			fmt.Printf("error removing %s \n", args)
		}
		return nil
	}

	// check for commands and execute them accordingly
	// check for help command
	// check for empty args or help flag
	if command.Flags().Changed("help") || len(args) < 1 {
		// run the help command
		help.Help(r, command, args)
		return nil
	}

	// check for interactive flag -i
	if r.Ii || command.Flags().Changed("interactive") {
		// run the interactive command
		interactive.InteractiveIi(args)
	}

	// check for interactive flag -I
	if r.II || command.Flags().Changed("INTERACTIVE") {
		// run the interactive command
		interactive.InteractiveII(args)
	}

	// check for the verbose flag
	if r.V || command.Flags().Changed("verbose") {
		// run the verbose command
		verbose.Verbose(args)
	}

	// remove empty directories
	if r.D || command.Flags().Changed("dir") {
		// remove empty directories
		directory.Directory(args)
	}

	// return on no error
	return nil
}

// ArgumentValidator - Validates the given arguments
func ArgumentValidator(rmx *rm.Rm) func(command *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		command_name := color.CyanString("rmx")
		// check if args are empty
		if len(args) < 1 {
			// print prompt
			fmt.Printf("%s: missing arguments or flags for command\nTry '%s --help' for more information.\n", command_name, command_name)
			return nil
		}

		// return on no error
		return nil
	}
}
