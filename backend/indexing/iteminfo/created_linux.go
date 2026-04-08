//go:build linux

package iteminfo

import (
	"os"
	"syscall"
	"time"
)

// GetCreatedTime returns the file's inode change time (ctime) on Linux.
// For newly created files this equals the creation time; it updates on renames/chmod too.
func GetCreatedTime(fi os.FileInfo) time.Time {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	}
	return fi.ModTime()
}
