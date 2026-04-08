//go:build !linux

package iteminfo

import (
	"os"
	"time"
)

// GetCreatedTime falls back to ModTime on non-Linux platforms.
func GetCreatedTime(fi os.FileInfo) time.Time {
	return fi.ModTime()
}
