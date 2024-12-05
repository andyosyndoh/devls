package internals

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func getLongFormat(path string, isDotEntry bool) string {
	var actualPath string

	// Clean the path first
	path = strings.TrimRight(path, "/")

	// Special handling for /dev directory and its parent
	if path == "/dev" {
		if isDotEntry {
			actualPath = "/dev" // For current directory
		} else {
			actualPath = "/" // For parent directory
		}
	} else if path == "/" {
		actualPath = "/"
	} else {
		actualPath = path
	}

	linkInfo, err := os.Lstat(actualPath)
	if err != nil {
		return ""
	}

	mode := linkInfo.Mode()
	nlink := linkInfo.Sys().(*syscall.Stat_t).Nlink
	uid := linkInfo.Sys().(*syscall.Stat_t).Uid
	gid := linkInfo.Sys().(*syscall.Stat_t).Gid
	size := linkInfo.Size()
	modTime := linkInfo.ModTime()

	var timeStr string
	if time.Since(modTime) < 6*30*24*time.Hour {
		timeStr = modTime.Format("Jan _2 15:04")
	} else {
		timeStr = modTime.Format("Jan _2  2006")
	}

	username := strconv.FormatUint(uint64(uid), 10)
	groupname := strconv.FormatUint(uint64(gid), 10)

	if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		username = u.Username
	}
	if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupname = g.Name
	}

	modeStr := mode.String()
	// if hasExtendedAttributes(actualPath) {
	// 	modeStr += "+"
	// }

	// Correctly set the first character for special files
	if mode&os.ModeSymlink != 0 {
		modeStr = "l" + modeStr[1:]
	} else if mode&os.ModeCharDevice != 0 {
		modeStr = "c" + modeStr[1:]
	} else if mode&os.ModeDevice != 0 {
		modeStr = "b" + modeStr[1:]
	} else if mode&os.ModeDir != 0 {
		modeStr = "d" + modeStr[1:]
	}

	// Handle sticky bit
	if mode&os.ModeSticky != 0 {
		if mode&0o002 != 0 {
			modeStr = modeStr[:len(modeStr)-1] + "t"
		} else {
			modeStr = modeStr[:len(modeStr)-1] + "T"
		}
	}

	// Remove duplicate 'c' if present
	if strings.HasPrefix(modeStr, "cc") {
		modeStr = "c" + modeStr[2:]
	}

	var result string
	if mode&os.ModeCharDevice != 0 || mode&os.ModeDevice != 0 {
		stat := getDeviceStat(actualPath)
		major, minor := majorMinor(stat.Rdev)
		result = fmt.Sprintf("%-10s %*d %-*s %-*s %*d, %*d %s",
			modeStr, LinkLen, nlink, UserLen, username, GroupLen, groupname,
			MajorLen, major, MinorLen, minor, timeStr)
	} else {
		result = fmt.Sprintf("%-10s %*d %-*s %-*s %*d %s",
			modeStr, LinkLen, nlink, UserLen, username, GroupLen, groupname,
			SizeLen, size, timeStr)
	}

	return result
}

// func hasExtendedAttributes(path string) bool {
// 	size, err := syscall.Listxattr(path, nil)
// 	if err != nil {
// 		return err != syscall.ENOTSUP
// 	}

// 	if size == 0 {
// 		return false
// 	}

// 	buf := make([]byte, size)
// 	_, err = syscall.Listxattr(path, buf)
// 	return err == nil
// }
