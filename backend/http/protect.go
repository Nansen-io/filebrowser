package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gtsteffaniak/filebrowser/backend/adapters/fs/files"
	"github.com/gtsteffaniak/filebrowser/backend/chainfs"
	"github.com/gtsteffaniak/filebrowser/backend/common/settings"
	"github.com/gtsteffaniak/filebrowser/backend/common/utils"
	"github.com/gtsteffaniak/filebrowser/backend/database/users"
	"github.com/gtsteffaniak/go-logger/logger"
	"golang.org/x/sys/unix"
)

const xattrFileGuid = "user.chainfs.fileguid"
const xattrProtectExpiry = "user.chainfs.protectexpiry"
const segmentThreshold = 10 * 1024 * 1024 // 10 MB

// protectHandler uploads a file to ChainFS and makes it read-only on disk.
// POST /api/chainfs/protect?path=<path>&source=<source>&hours=<hours>
func protectHandler(w http.ResponseWriter, r *http.Request, d *requestContext) (int, error) {
	encodedPath := r.URL.Query().Get("path")
	source := r.URL.Query().Get("source")
	hoursStr := r.URL.Query().Get("hours")

	filePath, err := url.QueryUnescape(encodedPath)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid path encoding: %w", err)
	}
	source, err = url.QueryUnescape(source)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid source encoding: %w", err)
	}

	hours := 24
	if hoursStr != "" {
		parsed, parseErr := strconv.Atoi(hoursStr)
		if parseErr != nil || parsed < 1 {
			return http.StatusBadRequest, fmt.Errorf("hours must be a positive integer")
		}
		hours = parsed
	}

	// Require ChainFS login
	if d.user.LoginMethod != users.LoginMethodChainFs || d.user.AzureAccessToken == "" {
		return http.StatusForbidden, fmt.Errorf("ChainFS account required to protect files")
	}

	// Check token expiry
	if d.user.AzureTokenExpiry > 0 && time.Now().Unix() > d.user.AzureTokenExpiry {
		return http.StatusUnauthorized, fmt.Errorf("ChainFS token expired, please re-authenticate")
	}

	// Check subscription
	if !d.user.ChainFSSubscribed {
		return http.StatusPaymentRequired, fmt.Errorf("an active ChainFS subscription is required to protect files")
	}

	// Resolve the real path on disk
	userScope, err := settings.GetScopeFromSourceName(d.user.Scopes, source)
	if err != nil {
		return http.StatusForbidden, err
	}
	userScope = strings.TrimRight(userScope, "/")

	fileInfo, err := files.FileInfoFaster(utils.FileOptions{
		Username:   d.user.Username,
		Path:       utils.JoinPathAsUnix(userScope, filePath),
		Source:     source,
		Expand:     false,
		ShowHidden: d.user.ShowHidden,
	}, store.Access)
	if err != nil {
		return errToStatus(err), err
	}
	if fileInfo.Type == "directory" {
		return http.StatusBadRequest, fmt.Errorf("cannot protect a directory")
	}

	// Decrypt the stored Azure token
	accessToken, err := decryptToken(d.user.AzureAccessToken)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// Derive per-user AES password
	aesPassword := deriveUserAESPassword(d.user)

	chainfsConfig := settings.Config.Auth.Methods.ChainFsAuth

	// Open the file
	f, err := os.Open(fileInfo.RealPath)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	stat, err := os.Stat(fileInfo.RealPath)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to stat file: %w", err)
	}

	// Upload to ChainFS (segmented if >10MB)
	var fileGuid string
	if stat.Size() > segmentThreshold {
		fileGuid, err = chainfs.UploadFileSegmented(chainfsConfig.ApiBaseUrl, accessToken, stat.Name(), f, stat.Size(), aesPassword)
	} else {
		fileGuid, err = chainfs.UploadFile(chainfsConfig.ApiBaseUrl, accessToken, stat.Name(), f, aesPassword)
	}
	if err != nil {
		logger.Errorf("ChainFS upload failed for %s: %v", fileInfo.RealPath, err)
		return http.StatusBadGateway, fmt.Errorf("ChainFS upload failed: %w", err)
	}
	logger.Infof("ChainFS upload succeeded for %s, FileGuid: %s", fileInfo.RealPath, fileGuid)

	// Store FileGuid in xattr
	if err := unix.Setxattr(fileInfo.RealPath, xattrFileGuid, []byte(fileGuid), 0); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to set xattr: %w", err)
	}

	// Store protection expiry as a little-endian int64 Unix timestamp
	expiry := time.Now().Add(time.Duration(hours) * time.Hour).Unix()
	expiryBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(expiryBuf, uint64(expiry))
	if err := unix.Setxattr(fileInfo.RealPath, xattrProtectExpiry, expiryBuf, 0); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to set expiry xattr: %w", err)
	}

	// Make the file read-only
	if err := os.Chmod(fileInfo.RealPath, 0444); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to set file read-only: %w", err)
	}

	return renderJSON(w, r, map[string]string{"fileGuid": fileGuid, "protectedUntil": time.Unix(expiry, 0).UTC().Format(time.RFC3339)})
}

// deriveUserAESPassword creates a stable per-user AES password from the server auth key + username.
func deriveUserAESPassword(user *users.User) string {
	material := settings.Config.Auth.Key + ":" + user.Username
	hash := sha256.Sum256([]byte(material))
	return hex.EncodeToString(hash[:])
}

// IsFileProtected returns true if the file at realPath has a ChainFS FileGuid xattr.
func IsFileProtected(realPath string) bool {
	buf := make([]byte, 64)
	n, err := unix.Getxattr(realPath, xattrFileGuid, buf)
	return err == nil && n > 0
}

// IsProtectionActive returns true if the file is protected AND its expiry has not yet passed.
// Files with no expiry xattr are treated as permanently protected.
func IsProtectionActive(realPath string) bool {
	if !IsFileProtected(realPath) {
		return false
	}
	buf := make([]byte, 8)
	n, err := unix.Getxattr(realPath, xattrProtectExpiry, buf)
	if err != nil || n != 8 {
		// No expiry stored — treat as permanently protected
		return true
	}
	expiry := int64(binary.LittleEndian.Uint64(buf))
	return time.Now().Unix() < expiry
}

// ProtectionExpiresAt returns the Unix timestamp when protection expires, and whether one is set.
func ProtectionExpiresAt(realPath string) (int64, bool) {
	buf := make([]byte, 8)
	n, err := unix.Getxattr(realPath, xattrProtectExpiry, buf)
	if err != nil || n != 8 {
		return 0, false
	}
	return int64(binary.LittleEndian.Uint64(buf)), true
}
