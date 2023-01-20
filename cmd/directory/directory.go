package directory

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/braswelljr/rmx/internal/err"
	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// Directory - checks if the given directory exists and is a directory and not a file and not empty and removes the directory
//
//	@param {[]string} args - the directory to remove
//	@return {error} - error if there is one
func Directory(args []string) {
	// check if the args are empty
	if len(args) < 1 {
		fmt.Print(err.ErrInvalidDirectory)
		return
	}

	cmd_name := color.CyanString("rmx")

	// loop through the directories
	for i, arg := range args {
		// check if the directory exists
		if exists, err := utils.Exists(arg); err != nil || !exists {
			if i == 0 {
				fmt.Printf("%s: cannot remove '%s': directory does not exist", cmd_name, arg)
			} else {
				fmt.Printf("\n%s: cannot remove '%s': directory does not exist", cmd_name, arg)
			}
			continue
		}

		// check if the directory is a directory
		if isDir, err := utils.IsDirectory(arg); err != nil || !isDir {
			if i == 0 {
				fmt.Printf("%s: '%s' is not a directory", cmd_name, arg)
			} else {
				fmt.Printf("\n%s: '%s' is not a directory", cmd_name, arg)
			}
			continue
		}

		// check if the directory is empty
		if isEmpty := utils.IsEmpty(arg); !isEmpty {
			if i == 0 {
				fmt.Printf("%s: cannot remove '%s': Directory not empty", cmd_name, arg)
			} else {
				fmt.Printf("\n%s: cannot remove '%s': Directory not empty", cmd_name, arg)
			}
			continue
		}

		// remove the directory
		if err := rm.RemoveDirectory(arg); err != nil {
			if i == 0 {
				fmt.Printf("%s: unable to remove directory '%s'", cmd_name, arg)
			} else {
				fmt.Printf("\n%s: unable to remove directory '%s'", cmd_name, arg)
			}
		}
	}
}
