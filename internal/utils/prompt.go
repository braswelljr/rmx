package utils

import (
	"fmt"

	"github.com/fatih/color"
)

// InteractivePrompt - the prompt to display to the user.
type InteractivePrompt struct {
	Prompt       string
	Type         string // warning, error, info, success, danger, default
	Confirmation bool
}

// Prompt - prompt the user for confirmation.
//
//	@param {string} prompt - the prompt to display to the user.
//	@param {string} color - the color of the prompt.
//	@param {string} confirmation - the confirmation string.
//	@return {bool} - true if the user confirmed the prompt.
func Prompt(prompt *InteractivePrompt) bool {
	// check for the prompt type to display the correct color
	switch prompt.Type {
	case "warning":
		color.New(color.FgYellow).Println(prompt.Prompt)

	case "danger":
		color.New(color.FgRed).Println(prompt.Prompt)

	case "error":
		color.New(color.FgRed).Println(prompt.Prompt)

	case "info":
		color.New(color.FgBlue).Println(prompt.Prompt)

	case "success":
		color.New(color.FgGreen).Println(prompt.Prompt)

	default:
		color.New(color.FgWhite).Println(prompt.Prompt)
	}

	// create the prompt
	fmt.Printf("Type %s for confirmation : ", color.MagentaString("[yes / y / Yes]"))

	// get the user confirmation
	var confirmation string

	// get the user confirmation
	fmt.Scanln(&confirmation)

	// check if the user confirmed the prompt
	if confirmation == "y" || confirmation == "Y" || confirmation == "yes" || confirmation == "Yes" {
		prompt.Confirmation = true
	} else {
		prompt.Confirmation = false
	}

	// check if the user confirmed the prompt
	return prompt.Confirmation
}
