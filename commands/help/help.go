package help

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/internal/execute"
	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// Help - show the help message.
//
//	@param {*rm.RM} r - contains the info for the rm command.
//	@param {*cobra.Command} commands - commands.
//	@param {[] string} args - command arguments.
//	@return void
func Help(_ *rm.RM, command *cobra.Command, args []string) {
	// show the help message
	// println("rmx - remove files or directories")
	// println("Usage: rmx [OPTION]... [FILE]...")
	// println("Remove (unlink) the FILE(s).")
	// println("  -f, --force           ignore nonexistent files and arguments, never prompt")
	// println("  -i                    prompt before every removal")
	// println("  -I                    prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	// println("  -r, -R                remove directories and their contents recursively")
	// println("  -d                    remove empty directories")
	// println("  -v                    explain what is being done")
	// println("      --help     display this help and exit")
	// println("      --version  output version information and exit")

	// check if the command is a root command
	if execute.IsRootCommand(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		// check if the command is hidden
		if command.Hidden {
			// print the usage of the command
			HelpUsage(command)
			return
		}

		// declare array of commands
		var coreCommands []string
		var additionalCommands []string
		// check for core commands and additonal commands
		for _, com := range command.Commands() {
			// check for hidden commands and short commands is not empty
			if com.Short != "" || com.Hidden {
				continue
			}

			c := utils.RPad(com.Name()+" : ", com.NamePadding()) + com.Short

			// check if the command is a core command
			if com.Annotations["core"] == "true" {
				coreCommands = append(coreCommands, c)
			} else {
				additionalCommands = append(additionalCommands, c)
			}
		}

		// If there are no core commands, assume everything is a core command
		if len(coreCommands) < 1 {
			coreCommands = additionalCommands
			additionalCommands = []string{}
		}

		// check for length of core commands
		if len(coreCommands) > 0 {
			command.Println("Core Commands:")
			command.Println(utils.Indent(utils.Dedent(strings.Join(coreCommands, "  ")), 2))

			// check for length of additional commands
			if len(additionalCommands) > 0 {
				command.Println("Additional Commands:")
				command.Println(utils.Indent(utils.Dedent(strings.Join(additionalCommands, "  ")), 2))
			}

			// print the usage of the command
			HelpUsage(command)
			return
		}
	}

	// print the help message
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
