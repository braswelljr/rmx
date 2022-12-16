package help

import (
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// Help - show the help message.
//
//	@param {*rm.RM} r - contains the info for the rm command.
//	@param {*cobra.Command} commands - commands.
//	@param {[] string} args - command arguments.
//	@return {Error} err - error to be returned if any.
func Help(_ *rm.RM, _ *cobra.Command, _ []string) error {
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

// HelpUsage - prints the usage of the command
//
//	@param {*cobra.Command} command - command to print the usage of.
//	@return {string} - usage of the command.
func HelpUsage(command *cobra.Command) error {
	subcommands := command.Commands()

	// check for length of subcommands
	if len(subcommands) > 0 {
		command.Print("\n\nAvailable commands:\n")
		// print the subcommands
		for _, c := range subcommands {
			// check for hidden commands
			if c.Hidden {
				continue
			}
			command.Printf("  %s\n", c.Name())
		}
		return nil
	}

	// print the flag usage information
	flagUsages := command.LocalFlags().FlagUsages()
	// check if usages are empty
	if flagUsages != "" {
		command.Println("\n\nFlags:")
		command.Print(utils.Indent(utils.Dedent(flagUsages), 2))
	}
	return nil
}
