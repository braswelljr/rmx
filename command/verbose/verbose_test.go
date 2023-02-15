package verbose

import (
	"os"
	"testing"
)

// TestVerbose - tests the verbose command.
//
//	@param t - testing object.
//	@return void
func TestVerbose(t *testing.T) {
	// create temp files and directories
	dir := os.TempDir()

	file := dir + "/test.txt"

	// remove the file if it exists
	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			t.Error("Removing file: ", err)
		}
	}

	// create a file in the directory
	if _, err := os.Create(file); err != nil {
		t.Error("Creating file: ", err)
	}

	// remove the directory if it exists
	if _, err := os.Stat(dir + "/test"); err == nil {
		if err := os.Remove(dir + "/test"); err != nil {
			t.Error("Removing directory: ", err)
		}
	}

	// create a directory in the directory
	if err := os.Mkdir(dir+"/test", 0755); err != nil {
		t.Error("Creating directory: ", err)
	}

	// set args as paths for the command
	args := []string{dir + "/test.txt", dir + "/test"}

	// call the verbose command
	Verbose(args)
}
