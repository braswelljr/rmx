.. _one-file-system:

one-file-system
---------------

Stay on one file system during recursive removal (--one-file-system).

Synopsis
~~~~~~~~


When removing a hierarchy recursively, --one-file-system skips any subdirectory that lives on a different file system than the argument the removal started from. It guards against accidentally descending into a mounted volume. On Windows, where device identity is unavailable, the flag is a no-op.

PERMISSIONS. The flag compares device IDs (via stat) and needs no extra permissions beyond those recursive removal already requires. Skipped mount points are reported on stderr and are not treated as errors.

USE CASES.
  - Clearing a directory that contains mounted volumes without deleting across the mount.
  - Safely removing a tree on the root file system while leaving mounted disks intact.

::

  one-file-system [DIR]... [flags]

Examples
~~~~~~~~

::

    # Remove a tree but skip file systems mounted beneath it
    rmx -r --one-file-system /data

Options
~~~~~~~

::

  -h, --help              help for one-file-system
      --one-file-system   when removing recursively, skip directories on a different file system

