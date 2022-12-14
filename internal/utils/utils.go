package utils

import "os"

// IsHidden returns true if the file/directory is hidden
// @param path - the path to the file/directory
// @param forceHidden - true if the file/directory is forced to be hidden
// @return bool - true if the file/directory is hidden
func IsHidden(path string, forceHidden bool) bool {
	// check if the file/directory is forced to be hidden
	if forceHidden {
		return true
	}

	// check if the file/directory is hidden by the OS (e.g. .git, .vscode)
	return int(path[0]) == DotChar
}

// IsEmpty returns true if the directory is empty
// @param path - the path to the directory
// @return bool - true if the directory is empty
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
