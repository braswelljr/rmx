// Package force implements rm's -f/--force flag: ignore nonexistent files and
// arguments, never prompt, and suppress the resulting errors.
package force

import (
	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/core"
	"github.com/braswelljr/rmx/internal/docmeta"
)

// Register installs the -f/--force flag on fs.
func Register(fs *pflag.FlagSet) {
	fs.BoolP("force", "f", false, "ignore nonexistent files and arguments, never prompt")
}

// Meta returns the documentation metadata for the force flag.
func Meta() docmeta.Meta {
	return docmeta.Meta{
		Name:  "force",
		Use:   "force [FILE]...",
		Short: "Ignore nonexistent files and arguments, and never prompt.",
		Long: "The -f/--force flag makes rmx ignore files that do not exist (exiting " +
			"successfully instead of reporting an error) and suppresses every confirmation " +
			"prompt, including the write-protected-file prompt. It is what turns rmx into an " +
			"unattended, script-safe delete. -f only affects prompting and missing-file errors; " +
			"it does not enable directory removal — combine it with -r for that.",
		Permissions: "-f skips the write-protected-file prompt, so read-only files are removed " +
			"without asking. Deletion is still governed by the parent directory's permissions, " +
			"not the file's own mode: if you lack write+execute on the containing directory the " +
			"removal fails with 'Permission denied' even under -f.",
		UseCases: []string{
			"Unattended cleanup in CI pipelines and shell scripts where prompts would hang.",
			"Deleting paths that may or may not exist without treating absence as an error.",
			"Removing read-only or write-protected files without per-file confirmation.",
		},
		Examples: []docmeta.Example{
			{
				Description: "Remove a file, succeeding silently if it is already gone:",
				Bash:        "rmx -f maybe-missing.log",
				PowerShell:  "rmx -f maybe-missing.log",
			},
			{
				Description: "Recursively delete a build directory with no prompts:",
				Bash:        "rmx -rf build/",
				PowerShell:  "rmx -rf build",
			},
		},
	}
}

// Apply reads the force flag and folds it into o. Force also drops the prompt
// mode to PromptNever; a later interactive flag may raise it again.
func Apply(fs *pflag.FlagSet, o *core.Options) {
	force, _ := fs.GetBool("force")
	o.Force = force
	if force {
		o.Interactive = core.PromptNever
	}
}
