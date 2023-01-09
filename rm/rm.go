package rm

import (
	"errors"
	"os"

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
	return nil
}
