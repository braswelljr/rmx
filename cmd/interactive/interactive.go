package interactive

import (
	"fmt"
	"sync"

	"github.com/fatih/color"

	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// InteractiveIi - prompt before every removal.
//
//	@param {[] string} args - command arguments.
//	@return void
func InteractiveIi(args []string) {
	// prompt the user for confirmation before removing the files
	if len(args) == 0 {
		args = []string{"."}
	}

	// walk through the files and directories and prompt the user for confirmation
	// before removing the files
	for _, arg := range args {
		// check if the arg is a directory
		if utils.IsDirectory(arg) {
			directory, err := utils.WalkDirectory(arg)
			if err != nil {
				continue
			}

			// work on arg as long as its a directory
			recursiveWorkOnDir(directory)
		} else {
			// prompt the user for confirmation
			prompt := &utils.InteractivePrompt{
				Prompt: fmt.Sprintf("\nAre you sure you want to remove the following file %s?", arg),
				Type:   "danger",
			}

			// prompt the user for confirmation
			if utils.Prompt(prompt) {
				// remove the files
				if err := rm.Remove(arg); err != nil {
					fmt.Printf("error removing %s \n", arg)
					continue
				}

				// print the file that was removed
				fmt.Printf("removed %s \n", color.RedString(arg))
			}
		}
	}
}

// recursiveWorkOnDir - recursively work on directory
//
//	@param directory - directory to be worked on
//	@return void
func recursiveWorkOnDir(directory *utils.Directory) {
	// Use a wait group to wait for all the directories/files to be traversed
	var wg sync.WaitGroup

	// recursively loop the directory.Files and directory.Directories.Files
	// and prompt the user for confirmation before removing the files
	for _, file := range directory.Files {
		// prompt the user for confirmation
		prompt := &utils.InteractivePrompt{
			Prompt: fmt.Sprintf("\nAre you sure you want to remove the following file %s?", file.Name),
			Type:   "danger",
		}

		// prompt the user for confirmation
		if utils.Prompt(prompt) {
			// remove the files
			if err := rm.Remove(file.Path); err != nil {
				fmt.Printf("error removing %s \n", file.Name)
				continue
			}

			// print the file that was removed
			fmt.Printf("removed %s \n", color.RedString(file.Name))
		}

	}

	// Use a channel to perform the directory traversal concurrently
	dirChan := make(chan *utils.Directory)

	// recursively loop the directory.Directories.Files and prompt the user for confirmation
	// before removing the files
	for _, directory := range directory.Directories {
		// Add to the wait group
		wg.Add(1)

		// Start a goroutine to walk the directory
		go func(dir *utils.Directory) {

			// Defer the wait group done
			defer wg.Done()

			if utils.IsEmpty(dir.Path) {
				// remove the directory
				if err := rm.Remove(dir.Path); err != nil {
					fmt.Printf("error removing %s \n", dir.Name)
				}
			}

			// Send the directory to the channel
			dirChan <- dir

			recursiveWorkOnDir(dir)
		}(directory)

	}

	// Close the channel when all the directories have been traversed
	go func() {
		wg.Wait()
		close(dirChan)
	}()
}

// InteractiveII - prompt once before removing all files.
//
//	@param {[] string} args - command arguments.
//	@return void
func InteractiveII(args []string) {
	// prompt the user for confirmation before removing the files
	if len(args) == 0 {
		args = []string{"."}
	}

	// create a new prompt
	var prompt *utils.InteractivePrompt

	if len(args) > 1 {
		// prompt the user for confirmation before removing the directory
		prompt = &utils.InteractivePrompt{
			Prompt: fmt.Sprintf("\nAre you sure you want to remove the following directories %s?", args),
			Type:   "danger",
		}
	} else {
		// prompt the user for confirmation before removing the directory
		prompt = &utils.InteractivePrompt{
			Prompt: fmt.Sprintf("\nAre you sure you want to remove the following directory %s?", args[0]),
			Type:   "danger",
		}
	}

	// prompt the user for confirmation
	if utils.Prompt(prompt) {
		// remove the files
		for _, arg := range args {
			if err := rm.RemoveAll(arg); err != nil {
				fmt.Printf("error removing %s \n", arg)
				continue
			}

			// print the file that was removed
			fmt.Printf("removed %s \n", color.RedString(arg))
		}
	}
}
