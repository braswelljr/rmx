.. _rmx:

rmx
---

Remove (unlink) files and directories — a cross-platform drop-in for rm.

Synopsis
~~~~~~~~


rmx removes each FILE. By default it does not remove directories; use -r to
remove a directory and everything under it, or -d to remove an empty directory.

::

  rmx [OPTION]... [FILE]... [flags]

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

