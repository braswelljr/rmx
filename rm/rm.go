package rm

import (
	"github.com/spf13/cobra"
)

// ShowHelp - show the help message
func (rm *RM) ShowHelp() error {
	// show the help message
	println("rmx - remove files or directories")
	println("Usage: rmx [OPTION]... [FILE]...")
	println("Remove (unlink) the FILE(s).")
	println("  -f, --force           ignore nonexistent files and arguments, never prompt")
	println("  -i                    prompt before every removal")
	println("  -I                    prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	println("  -r, -R                remove directories and their contents recursively")
	println("  -d                    remove empty directories")
	println("  -v                    explain what is being done")
	println("      --help     display this help and exit")
	println("      --version  output version information and exit")

	return nil
}

// Run - runs the command
// @param command - commands to run
// @param args - arguments and flags to run alongside the command
// @return error - error if there's one
func (rm *RM) Run(command *cobra.Command, args []string) error {
	// check for empty args or help flag
	if command.Flags().NFlag() < 1 || command.Flags().Changed("help") {
		return rm.ShowHelp()
	}

	// check for commands and execute them accordingly

	// return on no error
	return nil
}

// ArgumentValidator - Validates the given arguments
func ArgumentValidator(rm *RM) func(command *cobra.Command, args []string) error {
	return func(command *cobra.Command, args []string) error {
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
