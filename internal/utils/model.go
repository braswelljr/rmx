package utils

import "io/fs"

// Directory - a directory structure
type Directory struct {
	Name        string        `json:"name"`
	Files       []fs.DirEntry `json:"files"`
	Directories []*Directory  `json:"directories"`
}
