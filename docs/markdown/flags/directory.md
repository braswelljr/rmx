## directory

Remove empty directories (`-d`/`--dir`).

### Description

`-d`/`--dir` removes directories that are empty, the way `rmdir` does. It is a safe middle ground between plain `rmx` (which refuses directories entirely) and `-r` (which deletes contents too): a non-empty directory is reported as an error and left untouched.

### Permissions

Removing an empty directory requires write+execute permission on its parent. On a terminal, a write-protected directory prompts before removal unless `-f` is given. A non-empty directory reports 'Directory not empty' rather than a permission error.

### Use cases

- Pruning a single empty directory without risking its (future) contents.
- Removing leftover scaffold directories after their files were deleted.
- A safer alternative to `-r` when you expect the directory to be empty.

### Flags

```text
  -d, --dir   remove empty directories
```

### Examples

Remove a directory only if it contains no entries:

```bash
rmx -d generated/
```

```powershell
rmx -d generated
```
