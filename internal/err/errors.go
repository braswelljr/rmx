package err

import (
	"errors"
)

var (
	// ErrInvalidPath - error for invalid path
	ErrInvalidPath = errors.New("invalid path")
	// ErrInvalidPattern - error for invalid pattern
	ErrInvalidPattern = errors.New("invalid pattern")
	// ErrInvalidFile - error for invalid file
	ErrInvalidFile = errors.New("invalid file")
	// ErrInvalidDirectory - error for invalid directory
	ErrInvalidDirectory = errors.New("invalid directory")
	// ErrInvalidFileOrDirectory - error for invalid file or directory
	ErrInvalidFileOrDirectory = errors.New("invalid file or directory")
	// ErrEmptyDirectory - error for empty directory
	ErrEmptyDirectory = errors.New("empty directory")
	// ErrDirectoryDoesNotExist - error for directory does not exist
	ErrDirectoryDoesNotExist = errors.New("directory does not exist")
)
