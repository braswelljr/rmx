package utils

import "strings"

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
