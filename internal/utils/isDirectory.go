package utils

import "os"

// IsDirectory returns true if the path is a directory.
//
//	@param path - the path to the file/directory.
//	@return bool - true if the path is a directory.
func IsDirectory(path string) bool {
	// get the file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	// check if the path is a directory
	return fileInfo.IsDir()
}
