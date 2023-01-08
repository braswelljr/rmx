package interactive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// Interactive - prompt before every removal.
//
//	@param {*rm.RM} r - contains the info for the rm command.
//	@param {*cobra.Command} commands - commands.
//	@param {[] string} args - command arguments.
//	@return void
func Interactive(_ *rm.RM, _ *cobra.Command, args []string) {
	// prompt the user for confirmation before removing the files
	if len(args) == 0 {
		args = []string{"."}
	}

	// walk through the files and directories and prompt the user for confirmation
	// before removing the files
	for _, arg := range args {
		directory, err := utils.WalkDirectory(arg)
		if err != nil {
			continue
		}

		// recursively loop the directory.Files and directory.Directories.Files
		// and prompt the user for confirmation before removing the files
		for _, file := range directory.Files {
			// prompt the user for confirmation
			prompt := &utils.InteractivePrompt{
				Prompt: fmt.Sprintf("Are you sure you want to remove the following file %s?", file.Name()),
				Type:   "danger",
			}

			// get the file path and name
			filePath, _ := file.Info()

			// prompt the user for confirmation
			if utils.Prompt(prompt) {
				// remove the files
				fmt.Print(filePath)
			}

		}

		// recursively loop the directory.Directories.Files and prompt the user for confirmation
		// before removing the files
		for _, directory := range directory.Directories {
			for _, file := range directory.Files {
				// prompt the user for confirmation
				prompt := &utils.InteractivePrompt{
					Prompt: fmt.Sprintf("Are you sure you want to remove the following file %s?", file.Name()),
					Type:   "danger",
				}

				// get the file path and name
				filePath, _ := file.Info()

				// prompt the user for confirmation
				if utils.Prompt(prompt) {
					// remove the files
					fmt.Print(filePath)
				}
			}
		}
	}
}
