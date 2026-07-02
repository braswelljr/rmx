// Package version implements rm's --version flag.
package version

import (
	"fmt"
	"io"

	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/internal/common"
	"github.com/braswelljr/rmx/internal/docmeta"
)

// Meta returns the documentation metadata for the version flag.
func Meta() docmeta.Meta {
	return docmeta.Meta{
		Name:  "version",
		Use:   "version",
		Short: "Print version information and exit (--version).",
		Long:  "--version prints the rmx version string (injected at build time) and exits without removing anything.",
		Permissions: "Performs no filesystem access and needs no permissions; it never removes " +
			"anything, even when paths are also supplied.",
		UseCases: []string{
			"Confirming which rmx build is installed.",
			"Capturing the version in CI logs or bug reports.",
		},
		Examples: []docmeta.Example{
			{
				Description: "Print the rmx version:",
				Bash:        "rmx --version",
				PowerShell:  "rmx --version",
			},
		},
	}
}

// Register installs the --version flag on fs.
func Register(fs *pflag.FlagSet) {
	fs.Bool("version", false, "output version information and exit")
}

// Requested reports whether --version was given.
func Requested(fs *pflag.FlagSet) bool {
	v, _ := fs.GetBool("version")
	return v
}

// Print writes the version line to w.
func Print(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", common.AppName, common.Version)
}
