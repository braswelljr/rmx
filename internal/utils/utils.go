package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// IsDirectory returns true if the path is a directory.
//
//	@param path - the path to the file/directory.
//	@return bool - true if the path is a directory.
func IsDirectory(path string) (bool, error) {
	// get the file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("unable to get file info for %s: %w", path, err)
	}

	// check if the path is a directory
	return fileInfo.IsDir(), nil
}

// IsEmpty returns true if the directory is empty
//
//	@param path - the path to the directory
//	@return bool - true if the directory is empty
func IsEmpty(path string) bool {
	// check if the directory exists
	if exists, err := Exists(path); err != nil || !exists {
		return false
	}

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

// Exists - checks if a directory exists
//
//	@param {string} path - the path to the directory
//	@return {bool} - true if the directory exists
//	@return {error} - error
func Exists(path string) (bool, error) {
	// check if the directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}

// WalkDirectory - repeatedly steps through a directory and returns the
//
//	@param {string} root - the path to the directory
//	@return {*Directory} - the directory info
//	@return {error} - error
func WalkDirectory(root string) (*Directory, error) {
	// check if the directory exists
	if exists, err := Exists(root); err != nil || !exists {
		return &Directory{}, err
	}

	// for root and set it to current dir if its not set
	if len(strings.Trim(root, " ")) < 1 {
		root = "."
	}

	// open the directory
	dir, err := os.Open(root)
	if err != nil {
		return &Directory{}, err
	}

	// close file when not in use
	defer dir.Close()

	// read files in dir
	filesAndDirs, err := dir.ReadDir(-1)
	if err != nil {
		return &Directory{}, err
	}

	// create directory structure
	directory := &Directory{
		Name:        filepath.Base(root),
		Path:        root,
		IsEmpty:     IsEmpty(root),
		Files:       []*File{},
		Directories: []*Directory{},
	}

	// Use a channel to perform the directory traversal concurrently
	dirChan := make(chan *Directory)

	// Use a wait group to wait for all the directories to be traversed
	var wg sync.WaitGroup

	// iterate through the list of files and directories
	for _, fileOrDir := range filesAndDirs {
		// get the file path
		fpath := filepath.Join(root, fileOrDir.Name())

		// check if the file is a directory or a file
		if fileOrDir.IsDir() {
			// Add to the wait group
			wg.Add(1)

			// Start a goroutine to walk the directory
			go func(path string) {
				// Defer the wait group done
				defer wg.Done()

				// recursively walk directory
				subdir, err := WalkDirectory(path)
				if err != nil {
					return
				}

				// send the directory to the channel
				dirChan <- subdir

				// pass the directory path to the channel
			}(fpath)

			// receive the directory from the channel
			subdir := <-dirChan

			// append the directory to the list of directories
			directory.Directories = append(directory.Directories, subdir)

			// if the file is not a directory
		} else {
			info, _ := fileOrDir.Info()
			// create a file
			file := &File{
				Name: fileOrDir.Name(),
				Path: fpath,
				Info: info,
			}

			// append files
			directory.Files = append(directory.Files, file)
		}
	}

	// Close the channel when all the directories have been traversed
	go func() {
		wg.Wait()
		close(dirChan)
	}()

	// return the results
	return directory, nil
}

// RPad - adds padding to the right of a command
//
//	@param {string} s - string to be worked on
//	@param {int} padding - number of padding
//	@return string
func RPad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

// LPad - adds padding to the left of a command
//
//	@param {string} s - string to be worked on
//	@param {int} padding - number of padding
//	@return string
func LPad(s string, padding int) string {
	template := fmt.Sprintf("%%%ds ", padding)
	return fmt.Sprintf(template, s)
}

// Indent - indents the string
//
//	@param s - the string to indent
//	@param indent - the number of spaces to indent
//	@return string - the indented string
func Indent(s string, indent int) string {
	// create a string builder
	var builder strings.Builder

	// loop through the string
	for _, line := range strings.Split(s, " ") {
		// add the indent
		builder.WriteString(strings.Repeat(" ", indent))

		// add the line
		builder.WriteString(line)

		// add a new line
		builder.WriteString(" ")
	}

	// return the indented string
	return builder.String()
}

// Dedent - trims spaces on the left side of a line
//
//	@param {string} s - string to dedent
//	@return {string}
func Dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if l == "" {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

// GetFileExtension - gets the file extension
//
//	@param {string} filename - the filename
//	@return {string} - the file extension
func GetFileExtension(filename string) string {
	// get the file extension
	ext := filepath.Ext(filename)

	// remove the dot
	return strings.TrimPrefix(ext, ".")
}

// FileSize - gets the file size
//
//	@param {string} filename - the filename
//	@return {int64} - the file size
func FileSize(filename string) int64 {
	// get the file info
	info, _ := os.Stat(filename)

	// return the file size
	return info.Size()
}

// FileSizeToString - gets the file size and converts it to {bytes, KB, MB, GB}
//
//	@param {string} filename - the filename
//	@return {string} - the file size
func FileSizeToString(filename string) string {
	// get the file size
	size := FileSize(filename)

	// convert the file size to string
	return ByteCountSI(size)
}

// ByteCountSI - converts bytes to {bytes, KB, MB, GB}
//
//	@param {int64} b - the number of bytes
//	@return {string} - the converted bytes
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// GetFile - gets the file
//
//	@param {string} path - the path to the file
//	@return {File} - the file
func GetFile(path string) *File {
	// if the path is empty or the path is a directory
	if len(strings.TrimSpace(path)) == 0 {
		// return nil
		return nil
	}

	// check if the path is a directory and return nil if it is
	if isDir, err := IsDirectory(path); err != nil || isDir {
		// return nil
		return nil
	}

	// get the file info
	info, _ := os.Stat(path)

	// return the file
	return &File{
		Name: info.Name(),
		Path: path,
		Info: info,
	}
}

// GetFiles - gets the files in a directory
//
//	@param {string} path - the path to the directory
//	@return {[]File} - the files
func GetFiles(path string) []*File {
	// create a slice of files
	var files []*File

	// get the directory
	dir, _ := os.Open(path)

	// get the files
	fileInfos, _ := dir.Readdir(-1)

	// loop through the files
	for _, fileInfo := range fileInfos {
		// create a file
		file := &File{
			Name: fileInfo.Name(),
			Path: path + "/" + fileInfo.Name(),
			Info: fileInfo,
		}

		// append the file
		files = append(files, file)
	}

	// return the files
	return files
}
