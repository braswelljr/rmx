//go:build windows

package utils

import (
	"path/filepath"
	"syscall"

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

	// dotfiles also count as hidden (if you want)
	if int(path[0]) == constants.DotChar {
		return true
	}

	// get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	// Appending `\\?\` to the absolute path helps with
	// preventing 'Path Not Specified Error' when accessing
	// long paths and filenames
	// https://docs.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation?tabs=cmd
	pointer, err := syscall.UTF16PtrFromString(`\\?\` + absPath)
	if err != nil {
		return true
	}

	// get the file attribute if not specified in windows
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return true
	}

	// check for the file attribute
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
