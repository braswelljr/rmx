package utils

import (
	"os"

	"github.com/braswelljr/rmx/internal/constants"
)

// IsHidden returns true if the file/directory is hidden
//
//	@param path - the path to the file/directory
//	@param forceHidden - true if the file/directory is forced to be hidden
//	@return bool - true if the file/directory is hidden
func IsHidden(path string, forceHidden bool) bool {
	// check if the file/directory is forced to be hidden
	if forceHidden {
		return true
	}

	// check if the file/directory is hidden by the OS (e.g. .git, .vscode)
	return int(path[0]) == constants.DotChar
}

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

// GenerateRandomKey - generates a random key for the hill cipher
//
//	@param n - the size of the key
//	@return key - the key
func GenerateRandomKey(n int) string {
	// create a slice of bytes
	b := make([]byte, n)

	// loop through the key and add it to the key matrix
	for i, cache, remain := n-1, constants.SeedSrc.Int63(), constants.LetterIndexMask; i >= 0; {
		// if the cache is exhausted, get a new one
		if remain == 0 {
			cache, remain = constants.SeedSrc.Int63(), constants.LetterIndexMask
		} else {
			// get a random index from the key
			if idx := int(cache) & constants.LetterIndexMask; idx < len(constants.SourceKey) {
				// add the letter to the key
				b[i] = constants.SourceKey[idx]
				i--
			}
			// shift the cache and decrement the remain count
			cache >>= constants.LetterIndexBits
			remain--
		}

	}

	return string(b)
}
