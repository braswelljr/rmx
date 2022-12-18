package execute

import "github.com/spf13/cobra"

// IsRootCommand - check if the command is a root command
//
//	@param {*cobra.Command} command - command to check.
//	@return {bool} - true if the command is a root command, false if not.
func IsRootCommand(command *cobra.Command) bool {
	// check if the command is a root command
	return command != nil && !command.HasParent()
}
