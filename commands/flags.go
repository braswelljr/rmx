package commands

import (
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/rm"
)

// CommandFlags - the run command for flags
// @param command - commands to run
// @param args - arguments and flags to run alongside the command
func CommandFlags(command *cobra.Command, rmc *rm.RM) {
	command.PersistentFlags().BoolP("help", "h", false, "help for this command")
	command.Flags().BoolVarP(&rmc.Flags.F, "force", "f", false, "ignore nonexistent files and arguments, never prompt")
	command.Flags().BoolVarP(&rmc.Flags.Ii, "interactive", "i", false, "prompt before every removal")
	//command.Flags().BoolVarP(&rmc.Flags.II, "interactive", "I", false, "prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	command.Flags().BoolVarP(&rmc.Flags.Rr, "recursive", "r", false, "remove directories and their contents recursively")
	//command.Flags().BoolVarP(&rmc.Flags.RR, "recursive", "R", false, "remove directories and their contents recursively")
	command.Flags().BoolVarP(&rmc.Flags.D, "dir", "d", false, "remove empty directories")
	command.Flags().BoolVarP(&rmc.Flags.V, "verbose", "v", false, "explain what is being done")
	command.Flags().Float32Var(&rmc.Flags.Version, "version", 0.01, "output version information and exit")
}
