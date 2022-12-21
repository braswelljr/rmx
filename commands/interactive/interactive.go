package interactive

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/rm"
)

// Interactive - prompt before every removal.
//
//	@param {*rm.RM} r - contains the info for the rm command.
//	@param {*cobra.Command} commands - commands.
//	@param {[] string} args - command arguments.
//	@return void
func Interactive(r *rm.RM, command *cobra.Command, args []string) {
	// check if the flag is set for the command
	if r.Flags.Ii {
		// print the prompt
		fmt.Print(color.YellowString("rmx: remove '%s'? "))
		// check if the user wants to continue

	}
}
