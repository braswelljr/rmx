## onefilesystem

Stay on one file system during recursive removal (`--one-file-system`).

### Description

When removing a hierarchy recursively, `--one-file-system` skips any subdirectory that lives on a different file system than the argument the removal started from. It guards against accidentally descending into a mounted volume. On Windows, where device identity is unavailable, the flag is a no-op.

### Permissions

The flag compares device IDs (via stat) and needs no extra permissions beyond those recursive removal already requires. Skipped mount points are reported on stderr and are not treated as errors.

### Use cases

- Clearing a directory that contains mounted volumes without deleting across the mount.
- Safely removing a tree on the root file system while leaving mounted disks intact.

### Flags

```text
      --one-file-system   when removing recursively, skip directories on a different file system
```

### Examples

Remove a tree but skip file systems mounted beneath it:

```bash
rmx -r --one-file-system /data
```
