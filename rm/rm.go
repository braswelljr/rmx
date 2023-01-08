package rm

import (
	"os"
)

// Remove - remove the directory.
//
//	@param {string} path - the path to the directory.
//	@return {error} - error.
func Remove(path string) error {
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
	// remove the directory
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
