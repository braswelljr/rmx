## preserveroot

Refuse to recurse on '/' (`--preserve-root`, default; `--no-preserve-root` overrides).

### Description

`--preserve-root` is the default failsafe: `rmx` refuses to operate recursively on the filesystem root '/'. `--no-preserve-root` disables the failsafe for the rare, deliberate case where you really do mean to recurse from a root. Use it with extreme care.

### Permissions

This is a path check, not a permission check: `rmx` refuses '/' before it touches the filesystem, regardless of whether you would otherwise have permission. It never grants access you do not already have.

### Use cases

- Left at its default, protecting against a catastrophic '`rmx` `-r` /' typo.
- Deliberately overridden (`--no-preserve-root`) only in throwaway or container roots.

### Flags

```text
      --no-preserve-root   do not treat '/' specially
      --preserve-root      do not remove '/' (default) (default true)
```

### Examples

The default failsafe refuses to recurse on the root:

```bash
rmx -r /
```

```powershell
rmx -r C:\
```

Disable the failsafe (dangerous — only when you mean it):

```bash
rmx -r --no-preserve-root /mnt/scratch
```
