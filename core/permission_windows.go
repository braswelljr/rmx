//go:build windows

package core

import "os"

// writable reports whether path is writable. On Windows there is no access(2);
// approximate with the owner write bit / read-only attribute.
func writable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	return info.Mode().Perm()&0o200 != 0
}

// deviceID is unsupported on Windows, so --one-file-system is a no-op there.
func deviceID(_ os.FileInfo) (uint64, bool) {
	return 0, false
}
