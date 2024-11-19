package internals

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func getLongFormat(path string) string {
	//////fmt.Printf("Debug: Entering getLongFormat for path: %s\n", path)

	linkInfo, err := os.Lstat(path)
	if err != nil {
		// //fmt.Printf("Debug: Error with Lstat: %v\n", err)
		return ""
	}

	//////fmt.Printf("Debug: Lstat successful. File mode: %v\n", linkInfo.Mode())

	mode := linkInfo.Mode()
	nlink := linkInfo.Sys().(*syscall.Stat_t).Nlink
	uid := linkInfo.Sys().(*syscall.Stat_t).Uid
	gid := linkInfo.Sys().(*syscall.Stat_t).Gid
	size := linkInfo.Size()
	modTime := linkInfo.ModTime().Format("Jan  2 15:04")
	name := baseName(path)
	color := GetFileColor(linkInfo.Mode(), fmt.Sprint(linkInfo))

	// mt.Printf("Debug: Basic file info - Size: %d, ModTime: %s, Name: %s\n", size, modTime, name)

	username := strconv.FormatUint(uint64(uid), 10)
	groupname := strconv.FormatUint(uint64(gid), 10)

	if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		username = u.Username
	}
	if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupname = g.Name
	}

	//////fmt.Printf("Debug: Username: %s, Groupname: %s\n", username, groupname)

	modeStr := mode.String()
	colorlink := ""
	linked := ""

	if mode&os.ModeSymlink != 0 {
		// //fmt.Println("Debug: File is a symlink")
		modeStr = "l" + modeStr[1:] // Ensure the first character is 'l' for symlinks
		link, err := os.Readlink(path)
		if err == nil {
			linkInfo, err := os.Lstat(link)
			if err == nil { 
				colorlink = GetFileColor(linkInfo.Mode(), fmt.Sprint(linkInfo))
            }
			linked = fmt.Sprintf("-> %s%s%s", colorlink, link, Reset)
			//////fmt.Printf("Debug: Symlink target: %s\n", link)
			// For symlinks, we'll keep the size of the link itself, not the target
		}
	}

	result := fmt.Sprintf("%s %2d %s %s %6d %s %s%s%s %s", modeStr, nlink, username, groupname, size, modTime, color, name, Reset, linked)
	return result
}
