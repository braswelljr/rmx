## version

Print version information and exit (`--version`).

### Description

`--version` prints the `rmx` version string (injected at build time) and exits without removing anything.

### Permissions

Performs no filesystem access and needs no permissions; it never removes anything, even when paths are also supplied.

### Use cases

- Confirming which `rmx` build is installed.
- Capturing the version in CI logs or bug reports.

### Flags

```text
      --version   output version information and exit
```

### Examples

Print the `rmx` version:

```bash
rmx --version
```

```powershell
rmx --version
```
