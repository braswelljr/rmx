package color

import (
	"github.com/fatih/color"
)

// Color - a color
type Color struct {
	color.Attribute // color attribute
	color.Color     // color interface
}
