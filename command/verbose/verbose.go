package verbose

import (
	"fmt"
	"sync"

	"github.com/braswelljr/rmx/internal/utils"
	"github.com/braswelljr/rmx/rm"
)

// Verbose - verbose command
//
//	@param args - command arguments.
//	@return void
func Verbose(args []string) {
	// check for arguments
	if len(args) == 0 {
		args = []string{"."}
	}

	// show the files and directories that are being removed
	for _, arg := range args {
		// check if the arg is a directory
		isDir, err := utils.IsDirectory(arg)
		if err != nil {
			fmt.Printf("error checking '%s' is a directory \n", arg)
			continue
		}

		if isDir {
			// walk through the directory
			directory, err := utils.WalkDirectory(arg)
			// check for errors walking the directory and continue to the next arg
			if err != nil {
				continue
			}

			// repeatedly step through the directory and remove the files and directories
			recursiveWorkOnDir(directory)
		} else {
			// remove the files
			if err := rm.Remove(arg); err != nil {
				fmt.Printf("error removing '%s' \n", arg)
				continue
			} else {
				// get the file for the given path for file info
				file := utils.GetFile(arg)
				// check for errors getting the file

				if file == nil {
					// print the file that was removed
					fmt.Printf("error getting file info for '%s' \n", arg)
				} else {
					fmt.Printf("removed '%s' - (%s) \n", file.Path, utils.ByteCountSI(file.Info.Size()))
				}

			}
		}
	}
}

// recursiveWorkOnDir - recursively work on a directory.
//
//	@param directory - directory to work on.
//	@return void
func recursiveWorkOnDir(directory *utils.Directory) {
	// Use a wait group to wait for all the directories/files to be traversed
	var wg sync.WaitGroup

	// loop through the directories
	for _, file := range directory.Files {
		// remove the files
		if err := rm.Remove(file.Path); err != nil {
			fmt.Printf("error removing '%s' \n", file.Path)
			continue
		} else {
			// print the file that was removed
			fmt.Printf("removed '%s' - (%s) \n", file.Path, utils.ByteCountSI(file.Info.Size()))
		}
	}

	// // Use a channel to perform the directory traversal concurrently
	dirChan := make(chan *utils.Directory)

	// loop through the directories
	for _, dir := range directory.Directories {
		// increment the wait group
		wg.Add(1)

		// start a goroutine to walk the directory
		go func(dir *utils.Directory) {
			// defer the wait group done
			defer wg.Done()

			// check if the directory is empty
			if utils.IsEmpty(dir.Path) {
				// remove the directory
				if err := rm.Remove(dir.Path); err != nil {
					fmt.Printf("error removing directory '%s' \n", dir.Path)
				} else {
					// print the directory that was removed
					fmt.Printf("removed directory '%s' \n", dir.Path)
				}
			}

			// send the directory to the channel
			dirChan <- dir

			// recursively work on the directory
			recursiveWorkOnDir(dir)
		}(dir)
	}

	// Close the channel when all the directories have been traversed
	go func() {
		wg.Wait()
		close(dirChan)
	}()
}
