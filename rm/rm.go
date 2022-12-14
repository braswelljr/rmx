package rm

import (
	"github.com/spf13/cobra"
)

// ShowHelp - show the help message
func (rm *RM) ShowHelp() error {
	// show the help message
	println("rmx - remove files or directories")
	println("Usage: rmx [OPTION]... [FILE]...")
	println("Remove (unlink) the FILE(s).")
	println("  -f, --force           ignore nonexistent files and arguments, never prompt")
	println("  -i                    prompt before every removal")
	println("  -I                    prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	println("  -r, -R                remove directories and their contents recursively")
	println("  -d                    remove empty directories")
	println("  -v                    explain what is being done")
	println("      --help     display this help and exit")
	println("      --version  output version information and exit")
	println("Note:  rm does not remove directories.  Use rmdir to remove directories.")

	return nil
}

func (rm *RM) Run(command *cobra.Command, args []string) error {
	// check for empty args or help flag
	if command.Flags().NFlag() < 1 || command.Flags().Changed("help") {
		return rm.ShowHelp()
	}

	// close the directory
	return nil
}
