package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
)

// TestIsHidden - check if a file is hidden
func TestIsHidden(t *testing.T) {
	// create a temporary directory to test depending on the OS
	dir := os.TempDir()
	var filesAndDirs []string

	// read the directory and get all files/directories
	if err := filepath.Walk(dir, func(_ string, info os.FileInfo, _ error) error {
		// append the file/directory name to the slice
		filesAndDirs = append(filesAndDirs, info.Name())
		return nil
	}); err != nil {
		t.Error("Error reading directory", err)
	}

	// if there are no files/directories, create a hidden file and check if it is hidden
	if len(filesAndDirs) < 1 {
		// create two files, one hidden and one not hidden
		files := []string{"/test.txt", "/.test.txt"}
		// loop through the files and create them
		for _, file := range files {
			if _, err := os.Create(dir + file); err != nil {
				t.Error("Creating file: ", err)
			}
		}
	}

	// if there are files/directories, check if they are hidden
	t.Run("IsHidden", func(t *testing.T) {
		// check if the file is hidden
		if IsHidden(filesAndDirs[0], false) {
			assert.Equal(t, true, IsHidden(filesAndDirs[0], false), "File should be hidden")
		}
	})

	t.Run("IsNotHidden", func(t *testing.T) {
		// check if the file is not hidden
		if !IsHidden(filesAndDirs[1], false) {
			assert.Equal(t, false, IsHidden(filesAndDirs[1], false), "File should not be hidden")
		}
	})

}

// TestIsEmpty - check if a directory is empty
func TestIsEmpty(t *testing.T) {
	// isEmpty
	t.Run("IsEmpty", func(t *testing.T) {
		// create a temporary directory to test depending on the OS
		dir := os.TempDir()

		// check if the directory is empty
		if IsEmpty(dir) {
			assert.Equal(t, true, IsEmpty(dir), "Directory should be empty")
		}
	})

	// isNotEmpty
	t.Run("IsNotEmpty", func(t *testing.T) {
		// create a temporary directory to test depending on the OS
		dir := os.TempDir()

		// create a file in the directory
		if _, err := os.Create(dir + "/test.txt"); err != nil {
			t.Error("Creating file: ", err)
		}

		// check if the directory is empty
		if !IsEmpty(dir) {
			assert.Equal(t, false, IsEmpty(dir), "Directory should not be empty")
		}
	})
}
