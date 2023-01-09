package utils

import (
	"io/fs"
	"reflect"
)

// File - a file structure
type File struct {
	Name string      `json:"name"`
	Path string      `json:"path"`
	Info fs.FileInfo `json:"info"`
}

// Directory - a directory structure
type Directory struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	IsEmpty     bool         `json:"is_empty"`
	Files       []*File      `json:"files"`
	Directories []*Directory `json:"directories"`
}

// stringInSlice - check if a string is in a slice
//
//	@param {string} value - string to be checked
//	@param {[]string} slice - slice to be checked
//	@return {bool} - true if the string is in the slice
func stringInSlice(value string, slice []string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

// MergeStructs - merge multiple structs into one struct
//
//	@param {interface{}} structs - structs to be merged
//	@return interface{}
func MergeStructs(structs ...interface{}) reflect.Type {
	// create a slice of struct fields
	var structFields []reflect.StructField
	// create a slice of struct field names
	var structFieldNames []string

	// iterate through the structs
	for _, item := range structs {
		// get the type of the struct
		rt := reflect.TypeOf(item)

		// iterate through the fields of the struct
		for i := 0; i < rt.NumField(); i++ {
			// get the field
			field := rt.Field(i)

			// check if the field is already in the slice
			if !stringInSlice(field.Name, structFieldNames) {
				// append the field to the slice
				structFields = append(structFields, field)
				// append the field name to the slice
				structFieldNames = append(structFieldNames, field.Name)
			}
		}
	}

	// return the merged struct
	return reflect.StructOf(structFields)
}
