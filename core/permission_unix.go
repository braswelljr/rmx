//go:build !windows

package core

import (
	"os"
	"syscall"
)

// wOK is the write-permission bit for access(2).
const wOK = 0x2

// writable reports whether the calling process may write to path. rm uses this
// to decide whether a file is "write-protected" and warrants a prompt.
func writable(path string) bool {
	return syscall.Access(path, wOK) == nil
}

// deviceID returns the device number backing info, used by --one-file-system.
func deviceID(info os.FileInfo) (uint64, bool) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, false
	}
	return uint64(st.Dev), true
}
