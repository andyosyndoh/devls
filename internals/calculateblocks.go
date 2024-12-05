package internals

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func calculateTotalBlocks(path string, includeHidden bool) int64 {
	var total int64

	SizeLen, LinkLen, UserLen, GroupLen, MajorLen, MinorLen = 0, 0, 0, 0, 0, 0

	entries, err := os.ReadDir(path)
	if err != nil {
		return total
	}
	for i, names := range entries {
		name := names.Name()
		if i == 0 && includeHidden {
			total += CalculateParents(path)
		}
		if !includeHidden && isHidden(name) {
			continue
		}
		entry, err := os.Lstat(path + "/" + name)
		if err != nil {
			continue
		}
		blocksize := entry.Sys().(*syscall.Stat_t)
		total += blocksize.Blocks / 2
		size := strconv.FormatInt(entry.Size(), 10)
		linkCount := strconv.Itoa(int(blocksize.Nlink))
		SizeLen = max(SizeLen, len(size))
		LinkLen = max(LinkLen, len(linkCount))
		UserLen = max(UserLen, len(getUserName(int(blocksize.Uid))))
		GroupLen = max(GroupLen, len(getGroupName(int(blocksize.Gid))))

		// Check for device files and calculate their major/minor lengths
		if entry.Mode()&os.ModeCharDevice != 0 || entry.Mode()&os.ModeDevice != 0 {
			stat := getDeviceStat(path + "/" + name)
			major, minor := majorMinor(stat.Rdev)
			a := len(strconv.Itoa(int(major))) + len(strconv.Itoa(int(minor))) + 2
			if a >= SizeLen {
				SizeLen = a
				MajorLen = len(strconv.Itoa(int(major)))
				MinorLen = len(strconv.Itoa(int(minor)))
			}
		}

	}
	return total
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isHidden(name string) bool {
	return name[0] == '.'
}

func getGroupName(gid int) string {
	g, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		return strconv.Itoa(gid)
	}
	return g.Name
}

func getUserName(uid int) string {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return strconv.Itoa(uid)
	}
	return u.Username
}

func getDeviceStat(filePath string) *syscall.Stat_t {
	stat := &syscall.Stat_t{}
	err := syscall.Lstat(filePath, stat)
	if err != nil {
		fmt.Printf("Error getting device stat: %v\n", err)
	}
	return stat
}

func majorMinor(rdev uint64) (uint64, uint64) {
	major := (rdev >> 8) & 0xfff
	minor := (rdev & 0xff) | ((rdev >> 12) & 0xfff00)
	return major, minor
}

func CalculateParents(path string) int64 {
	var total int64
	paths := []string{".", ".."}

	for _, p := range paths {
		entry, err := os.Lstat(path + "/" + p)
		if err != nil {
			continue
		}
		blocksize := entry.Sys().(*syscall.Stat_t)
		total += blocksize.Blocks / 2
		size := strconv.FormatInt(entry.Size(), 10)
		linkCount := strconv.Itoa(int(blocksize.Nlink))
		SizeLen = max(SizeLen, len(size))
		LinkLen = max(LinkLen, len(linkCount))
		UserLen = max(UserLen, len(getUserName(int(blocksize.Uid))))
		GroupLen = max(GroupLen, len(getGroupName(int(blocksize.Gid))))

		// Check for device files and calculate their major/minor lengths
		if entry.Mode()&os.ModeCharDevice != 0 || entry.Mode()&os.ModeDevice != 0 {
			stat := getDeviceStat(path + "/" + p)
			major, minor := majorMinor(stat.Rdev)
			a := len(strconv.Itoa(int(major))) + len(strconv.Itoa(int(minor))) + 2
			if a >= SizeLen {
				SizeLen = a
				MajorLen = len(strconv.Itoa(int(major)))
				MinorLen = len(strconv.Itoa(int(minor)))
			}
		}
	}
	return total
}
