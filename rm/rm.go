package rm

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/braswelljr/rmx/internal/utils"
)

// Remove - remove the directory.
//
//	@param {string} path - the path to the directory.
//	@return {error} - error.
func Remove(path string) error {
	// check if the directory exists
	if exists, err := utils.Exists(path); err != nil || !exists {
		return errors.New("Directory does not exist")
	}

	// remove the directory
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

// RemoveAll - remove the directory and all of its contents.
//
//	@param {string} path - the path to the directory.
//	@return {error} - error.
func RemoveAll(path string) error {
	// check if the directory exists
	if exists, err := utils.Exists(path); err != nil || !exists {
		return errors.New("Directory does not exist")
	}

	// remove the directory
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}

// RemoveMultiple - remove multiple directories and all of their contents.
//
//	@param {[]string} paths - the paths to the directories.
//	@return {error} - error.
func RemoveMultiple(paths []string) error {
	var err error
	// create a wait group
	var wg sync.WaitGroup
	// make an error channel
	errChan := make(chan error, len(paths))
	// loop through the paths concurrently
	for _, path := range paths {
		// increment the wait group
		wg.Add(1)

		// remove the directory concurrently
		go func(path string) {
			// decrement the wait group
			defer wg.Done()

			// check if the directory exists
			if exists, err := utils.Exists(path); err != nil || !exists {
				if err != nil {
					// set the error
					fmt.Println("directory does not exist: ", path)
				}
			}

			// remove the directory
			if err := RemoveAll(path); err != nil {
				// set the error
				errChan <- err
			}
		}(path)
	}

	errChanLen := len(errChan)

	// check if there are any errors
	if errChanLen > 0 {
		// loop through the errors
		for i := 0; i < errChanLen; i++ {
			// set the error
			err = <-errChan
		}
	}

	// wait for all the directories to be removed
	wg.Wait()

	return err
}

// RemoveDirectory - remove the directory.
//
//	@param {string} path - the path to the directory.
//	@return {error} - error.
func RemoveDirectory(path string) error {
	// check if the directory exists
	if exists, err := utils.Exists(path); err != nil || !exists {
		return errors.New("Directory does not exist")
	}

	//check if is a directory
	if isDir, err := utils.IsDirectory(path); err != nil || !isDir {
		return errors.New("path is not a directory")
	}

	// remove the directory
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
