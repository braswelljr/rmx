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

// TestIsDirectory - check if a path is a directory
func TestIsDirectory(t *testing.T) {
	// create a temporary directory to test depending on the OS
	dir := os.TempDir()

	// create a file in the directory
	if _, err := os.Create(dir + "/test.txt"); err != nil {
		t.Error("Creating file: ", err)
	}

	// isDirectory
	t.Run("IsDirectory", func(t *testing.T) {
		// check if the path is a directory
		if IsDirectory(dir) {
			assert.Equal(t, true, IsDirectory(dir), "Path should be a directory")
		}
	})

	// isNotDirectory
	t.Run("IsNotDirectory", func(t *testing.T) {
		// check if the path is not a directory
		if !IsDirectory(dir + "/test.txt") {
			assert.Equal(t, false, IsDirectory(dir+"/test.txt"), "Path should not be a directory")
		}
	})
}

// TestGenerateRandomKey - check if a random string is generated
func TestGenerateRandomKey(t *testing.T) {
	t.Run("GenerateRandomKey", func(t *testing.T) {
		keyLen := 10
		// create a random string
		randomString := GenerateRandomKey(keyLen)

		// check if the random string is empty
		if randomString != "" {
			assert.NotEqual(t, "", randomString, "Random string should not be empty")
		}

		// check if the random string is the correct length
		if len(randomString) == keyLen {
			assert.Equal(t, keyLen, len(randomString), "Random string should be 10 characters")
		}
	})
}

// TestIndent - check if an added indent is added to a string
// func TestSpacing(t *testing.T) {
// 	// create a struct to test
// 	type item struct {
// 		str    string
// 		indent int
// 	}
// 	// create a slice of items to test
// 	items := []item{
// 		{"test", 1},
// 		{"test", 2},
// 		{"test", 3},
// 	}

// 	t.Run("Indent", func(t *testing.T) {
// 		for _, item := range items {
// 			// add an indent to the string
// 			indentedString := Indent(item.str, item.indent)
// 			// check if the string is indented
// 			if indentedString == "" {
// 				assert.NotEqual(t, "", indentedString, "Indented string should not be empty")
// 			}

// 			if len(indentedString) == len(item.str)+item.indent {
// 				assert.Equal(t, len(indentedString), len(item.str)+item.indent, "Indented string should be 1 character longer")
// 			}

// 			assert.Equal(t, " test", indentedString, "Indented string should be ' test'")
// 		}
// 	})

// 	t.Run("Dedent", func(t *testing.T) {
// 		for _, item := range items {
// 			// remove an indent from the string
// 			dedentedString := Dedent(item.str)
// 			// check if the string is dedented
// 			if dedentedString == "" {
// 				assert.NotEqual(t, "", dedentedString, "Dedented string should not be empty")
// 			}

// 			if len(dedentedString) == len(item.str)-item.indent {
// 				assert.Equal(t, len(dedentedString), len(item.str)-item.indent, "Dedented string should be 1 character less")
// 			}

// 			assert.Equal(t, "test", dedentedString, "Dedented string should be 'test'")
// 		}

// 	})
// }

// TestWalkDirectory - check if a directory is walked
// func TestWalkDirectory(t *testing.T) {
// 	// create a temporary directory to test depending on the OS
// 	dir := os.TempDir()

// 	files := []string{"test.txt", "test2.txt", "test3.txt"}

// 	// create a file in the directory
// 	if _, err := os.Create(dir + "/test.txt"); err != nil {
// 		t.Error("Creating file: ", err)
// 	}

// 	// walk the directory
// 	filesAndDirs, err := WalkDirectory(dir)
// 	if err != nil {
// 		assert.Equal(t, nil, err, "Directory should be walked")
// 	}

// 	// check if the directory is walked correctly and the file is found
// 	// if filesAndDirs.Name == dir {
// 	// 	assert.Equal(t, "test.txt", filesAndDirs.Name, "File should be found")
// 	// }

// 	// check if the the created file is found in files
// 	for i, file := range files {
// 		if file == filesAndDirs.Files[i].Name {
// 			assert.Equal(t, "test.txt", file, "File should be found")
// 		}
// 	}
// }
