package docs

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/braswelljr/rmx/internal/common"
)

// guideLong is the prose for the installation/usage guide. It avoids markdown
// headings so it renders cleanly across markdown, man, yaml and reST.
const guideLong = `rmx is a cross-platform, drop-in replacement for the GNU/UNIX rm command. It ` +
	`mirrors rm's flags and behavior — prompting, recursion, write-protected handling and the ` +
	`--preserve-root failsafe — while adding concurrent removal on the non-interactive path.

INSTALLATION. Install the latest version with go install, or build from source with make build ` +
	`(the binary lands in ./bin) or make install. The version string is injected at build time ` +
	`and reported by --version; the examples below show the exact commands.

USAGE. Invoke rmx exactly like rm. By default rmx removes files but refuses directories; pass ` +
	`-r (or -R) to remove a directory and its contents recursively, or -d to remove a single ` +
	`empty directory. Prompting is controlled by -i (before every removal), -I (once before a ` +
	`bulk or recursive removal) and --interactive; -f disables all prompts and ignores missing ` +
	`files, which is what makes it safe in scripts, and -v reports each removal.

SAFETY. rmx refuses to remove the "." and ".." entries and, by default, refuses to recurse on ` +
	`the filesystem root (override with --no-preserve-root). On a terminal, a write-protected ` +
	`file triggers a confirmation prompt unless -f is given. --one-file-system keeps a recursive ` +
	`removal from crossing into a mounted volume.

USE CASES. Everyday file cleanup; unattended deletion in CI or scripts with -f; cautious, ` +
	`per-file review with -i; a single confirmation before a large delete with -I; pruning an ` +
	`empty directory with -d; staying on one file system with --one-file-system; and auditing a ` +
	`recursive delete with -rv.`

// guideExample is the worked-examples block for the guide.
const guideExample = `  # install the latest rmx
  go install github.com/braswelljr/rmx@latest

  # remove a few files
  rmx notes.txt draft.md

  # recursively remove a directory, reporting each deletion
  rmx -rv build/

  # unattended cleanup: no prompts, ignore missing paths
  rmx -rf /tmp/cache

  # confirm every removal
  rmx -i important/*.conf

  # remove an empty directory only
  rmx -d generated/`

// guideCommand builds the synthetic command backing the guide page. It carries
// the full flag set so the guide's options section is a complete reference.
func guideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     common.AppName + " [OPTION]... [FILE]...",
		Short:   "Installation and usage guide for rmx, the drop-in rm replacement.",
		Long:    guideLong,
		Example: guideExample,
		Run:     func(*cobra.Command, []string) {},
	}
	c.DisableAutoGenTag = true
	for _, t := range flagTopics() {
		t.register(c.Flags())
	}
	return c
}

// GenerateGuide writes the installation/usage guide (guide.*) into each format
// subdirectory under base.
func GenerateGuide(base string) error {
	c := guideCommand()

	if err := writeFile(base, "markdown", "guide.md", func(w io.Writer) error {
		return genMarkdownTagged(c, w)
	}); err != nil {
		return err
	}
	if err := writeFile(base, "rest", "guide.rst", func(w io.Writer) error {
		return doc.GenReST(c, w)
	}); err != nil {
		return err
	}
	if err := writeFile(base, "yaml", "guide.yaml", func(w io.Writer) error {
		return doc.GenYaml(c, w)
	}); err != nil {
		return err
	}
	header := manHeader(common.AppName + "-guide")
	return writeFile(base, "man", "guide.1", func(w io.Writer) error {
		return doc.GenMan(c, header, w)
	})
}
