package utils

import "fmt"

// RPad - adds padding to the right of a command
//
//	@param {string} s - string to be worked on
//	@param {int} padding - number of padding
//	@return string
func RPad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}
