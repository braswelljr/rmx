package force

import (
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/rm"
)

type ForceC struct {
	Force   bool
	Process rm.Rm
}

// Force - force a command to run.
//
//	@param {*rm.RM} r - contains the info for the rm command.
//	@param {*cobra.Command} commands - commands.
//	@param {[] string} args - command arguments.
//	@return void
func Force(r *rm.Rm, command *cobra.Command, args []string) {
	// force the command to run
}
