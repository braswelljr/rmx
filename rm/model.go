package rm

import (
	"io"
	"io/fs"

	"github.com/fatih/color"
)

// Flags is a struct that contains all the flags for the rm command
type Flags struct { // rm -h
	F       bool    `flag:"force" short:"f" description:"ignore nonexistent files and arguments, never prompt"`                                                                                                               // rm -f
	Ii      bool    `flag:"interactive" short:"i" description:"prompt before every removal"`                                                                                                                                  // rm -i
	II      bool    `flag:"interactive" short:"I" description:"prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes"` // rm -I
	Rr      bool    `flag:"recursive" short:"r" description:"remove directories and their contents recursively"`                                                                                                              // rm -r
	RR      bool    `flag:"recursive" short:"R" description:"remove directories and their contents recursively"`                                                                                                              // rm -R
	D       bool    `flag:"dir" short:"d" description:"remove empty directories"`                                                                                                                                             // rm -d
	V       bool    `flag:"verbose" short:"v" description:"explain what is being done"`                                                                                                                                       // rm -v
	Version float32 `flag:"version" description:"output version information and exit"`                                                                                                                                        // rm --version
}

// RM is a struct that contains all the information for the rm command
type RM struct {
	Flags Flags

	Stdin  io.Reader
	Stdout io.Writer
	Color  color.Color

	Directory string
}

// Directory is a struct that contains all the information for a directory
type Directory struct {
	Name string
	Path string
	Info fs.FileInfo
}

// Directories is a slice of Directory
type Directories []Directory
