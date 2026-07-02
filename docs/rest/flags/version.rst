.. _version:

version
-------

Print version information and exit (--version).

Synopsis
~~~~~~~~


--version prints the rmx version string (injected at build time) and exits without removing anything.

PERMISSIONS. Performs no filesystem access and needs no permissions; it never removes anything, even when paths are also supplied.

USE CASES.
  - Confirming which rmx build is installed.
  - Capturing the version in CI logs or bug reports.

::

  version [flags]

Examples
~~~~~~~~

::

    # Print the rmx version
    rmx --version

Options
~~~~~~~

::

  -h, --help      help for version
      --version   output version information and exit

