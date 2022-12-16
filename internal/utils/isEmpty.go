package utils

import (
	"os"
)

// IsEmpty returns true if the directory is empty
//
//	@param path - the path to the directory
//	@return bool - true if the directory is empty
func IsEmpty(path string) bool {
	// get the directory
	dir, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// get the directory entries
	dirEntries, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// check if the directory is empty
	return len(dirEntries) < 1
}
