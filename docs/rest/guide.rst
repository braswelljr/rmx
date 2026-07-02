.. _rmx:

rmx
---

Installation and usage guide for rmx, the drop-in rm replacement.

Synopsis
~~~~~~~~


rmx is a cross-platform, drop-in replacement for the GNU/UNIX rm command. It mirrors rm's flags and behavior — prompting, recursion, write-protected handling and the --preserve-root failsafe — while adding concurrent removal on the non-interactive path.

INSTALLATION. Install the latest version with go install, or build from source with make build (the binary lands in ./bin) or make install. The version string is injected at build time and reported by --version; the examples below show the exact commands.

USAGE. Invoke rmx exactly like rm. By default rmx removes files but refuses directories; pass -r (or -R) to remove a directory and its contents recursively, or -d to remove a single empty directory. Prompting is controlled by -i (before every removal), -I (once before a bulk or recursive removal) and --interactive; -f disables all prompts and ignores missing files, which is what makes it safe in scripts, and -v reports each removal.

SAFETY. rmx refuses to remove the "." and ".." entries and, by default, refuses to recurse on the filesystem root (override with --no-preserve-root). On a terminal, a write-protected file triggers a confirmation prompt unless -f is given. --one-file-system keeps a recursive removal from crossing into a mounted volume.

USE CASES. Everyday file cleanup; unattended deletion in CI or scripts with -f; cautious, per-file review with -i; a single confirmation before a large delete with -I; pruning an empty directory with -d; staying on one file system with --one-file-system; and auditing a recursive delete with -rv.

::

  rmx [OPTION]... [FILE]... [flags]

Examples
~~~~~~~~

::

    # install the latest rmx
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
    rmx -d generated/

Options
~~~~~~~

::

  -d, --dir                             remove empty directories
  -f, --force                           ignore nonexistent files and arguments, never prompt
  -h, --help                            help for rmx
      --interactive string[="always"]   prompt according to WHEN: never, once, or always
      --no-preserve-root                do not treat '/' specially
      --one-file-system                 when removing recursively, skip directories on a different file system
      --preserve-root                   do not remove '/' (default) (default true)
  -i, --prompt-each                     prompt before every removal
  -I, --prompt-once                     prompt once before removing more than three files, or when removing recursively
  -r, --recursive                       remove directories and their contents recursively
  -v, --verbose                         explain what is being done
      --version                         output version information and exit

