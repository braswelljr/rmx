// Package verbose implements rm's -v/--verbose flag: report each removal.
package verbose

import (
	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/core"
	"github.com/braswelljr/rmx/internal/docmeta"
)

// Meta returns the documentation metadata for the verbose flag.
func Meta() docmeta.Meta {
	return docmeta.Meta{
		Name:  "verbose",
		Use:   "verbose [FILE]...",
		Short: "Explain what is being removed (-v/--verbose).",
		Long: "-v/--verbose prints a line for each item removed, e.g. removed 'notes.txt'. To " +
			"keep the report in the order you listed the arguments, verbose removal runs " +
			"sequentially rather than concurrently.",
		Permissions: "-v does not change permission handling or what may be removed; it only " +
			"reports the items that were successfully removed. Files it could not remove are " +
			"reported as errors on stderr as usual.",
		UseCases: []string{
			"Auditing exactly what a recursive or wildcard delete removed.",
			"Confirming a script deleted the files you expected.",
			"Troubleshooting by pairing -v with -i to see and approve each removal.",
		},
		Examples: []docmeta.Example{
			{
				Description: "Remove several files, printing a line per file:",
				Bash:        "rmx -v a.txt b.txt",
				PowerShell:  "rmx -v a.txt b.txt",
			},
			{
				Description: "Recursively remove a directory, reporting each removal:",
				Bash:        "rmx -rv dist/",
				PowerShell:  "rmx -rv dist",
			},
		},
	}
}

// Register installs the -v/--verbose flag on fs.
func Register(fs *pflag.FlagSet) {
	fs.BoolP("verbose", "v", false, "explain what is being done")
}

// Apply sets o.Verbose from the -v flag.
func Apply(fs *pflag.FlagSet, o *core.Options) {
	o.Verbose, _ = fs.GetBool("verbose")
}
