package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestIsHidden - check if a file is hidden
func TestIsHidden(t *testing.T) {
	t.Run("IsHidden", func(t *testing.T) {
		// create a temporary directory to test depending on the OS
		dir, err := ioutil.TempDir("", "test")
		if err != nil {
			t.Error("Creating directory: ", err)
		}
		var filesAndDirs []string
		if err := filepath.Walk(dir, func(_ string, info os.FileInfo, _ error) error {
			filesAndDirs = append(filesAndDirs, info.Name())
			return nil
		}); err != nil {
			t.Error("Error reading directory", err)
		}

		// if there are no files/directories, create a hidden file and check if it is hidden
		if len(filesAndDirs) < 1 {
			// create two files, one hidden and one not hidden
			if err := ioutil.WriteFile(dir+"/test.txt", []byte("test"), 0644); err != nil {
				t.Error("Creating file: ", err)
			}
			if err := ioutil.WriteFile(dir+"/.test.txt", []byte("test"), 0644); err != nil {
				t.Error("Creating file: ", err)
			}
		} else {
			// if there are files/directories, check if they are hidden
			for _, file := range filesAndDirs {
				if IsHidden(file, false) {
					t.Error("File should not be hidden")
				}
			}
		}
	})
}

// TestIsEmpty - check if a directory is empty
func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		// create a temporary directory to test depending on the OS
		dir, err := ioutil.TempDir("", "test")
		if err != nil {
			t.Error("Creating directory: ", err)
		}

		// check if the directory is empty
		if !IsEmpty(dir) {
			t.Error("Directory should be empty")
		}

		// create a file in the directory and check if it is empty
		if err := ioutil.WriteFile(dir+"/test.txt", []byte("test"), 0644); err != nil {
			t.Error("Creating file: ", err)
		}

		if IsEmpty(dir) {
			t.Error("Directory should not be empty")
		}
	})
}
