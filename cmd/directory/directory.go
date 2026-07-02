// Package directory implements rm's -d/--dir flag: remove empty directories.
package directory

import (
	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/core"
	"github.com/braswelljr/rmx/internal/docmeta"
)

// Meta returns the documentation metadata for the dir flag.
func Meta() docmeta.Meta {
	return docmeta.Meta{
		Name:  "directory",
		Use:   "dir [DIR]...",
		Short: "Remove empty directories (-d/--dir).",
		Long: "-d/--dir removes directories that are empty, the way rmdir does. It is a safe " +
			"middle ground between plain rmx (which refuses directories entirely) and -r " +
			"(which deletes contents too): a non-empty directory is reported as an error and " +
			"left untouched.",
		Permissions: "Removing an empty directory requires write+execute permission on its " +
			"parent. On a terminal, a write-protected directory prompts before removal unless -f " +
			"is given. A non-empty directory reports 'Directory not empty' rather than a " +
			"permission error.",
		UseCases: []string{
			"Pruning a single empty directory without risking its (future) contents.",
			"Removing leftover scaffold directories after their files were deleted.",
			"A safer alternative to -r when you expect the directory to be empty.",
		},
		Examples: []docmeta.Example{
			{
				Description: "Remove a directory only if it contains no entries:",
				Bash:        "rmx -d generated/",
				PowerShell:  "rmx -d generated",
			},
		},
	}
}

// Register installs the -d/--dir flag on fs.
func Register(fs *pflag.FlagSet) {
	fs.BoolP("dir", "d", false, "remove empty directories")
}

// Apply sets o.Dir from the -d flag.
func Apply(fs *pflag.FlagSet, o *core.Options) {
	o.Dir, _ = fs.GetBool("dir")
}
