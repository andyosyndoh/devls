package internals

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func getLongFormat(path string) string {
    linkInfo, err := os.Lstat(path)
    if err != nil {
        return ""
    }

    mode := linkInfo.Mode()
    nlink := linkInfo.Sys().(*syscall.Stat_t).Nlink
    uid := linkInfo.Sys().(*syscall.Stat_t).Uid
    gid := linkInfo.Sys().(*syscall.Stat_t).Gid
    size := linkInfo.Size()
    modTime := linkInfo.ModTime().Format("Jan 2 15:04")
    name := baseName(path)
    color := GetFileColor(linkInfo.Mode(), fmt.Sprint(linkInfo))

    username := strconv.FormatUint(uint64(uid), 10)
    groupname := strconv.FormatUint(uint64(gid), 10)

    if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
        username = u.Username
    }
    if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
        groupname = g.Name
    }

    modeStr := mode.String()
    colorlink := ""
    linked := ""

    if mode&os.ModeSymlink != 0 {
        modeStr = "l" + modeStr[1:]
        link, err := os.Readlink(path)
        if err == nil {
            linkInfo, err := os.Lstat(link)
            if err == nil { 
                colorlink = GetFileColor(linkInfo.Mode(), fmt.Sprint(linkInfo))
            }
            linked = fmt.Sprintf("-> %s%s%s", colorlink, link, Reset)
        }
    }

    result := ""
    if linkInfo.Mode()&os.ModeCharDevice != 0 || linkInfo.Mode()&os.ModeDevice != 0 {
        stat := getDeviceStat(path)
        major, minor := majorMinor(stat.Rdev)
        result = fmt.Sprintf("%-10s %*d %-*s %-*s %*d, %*d %s %s %s",
            modeStr[1:], LinkLen, nlink, UserLen, username, GroupLen, groupname,
            MajorLen, major, MinorLen, minor, modTime, color+name+Reset, linked)
    } else {
        result = fmt.Sprintf("%-10s %*d %-*s %-*s %*d %s %s %s",
            modeStr, LinkLen, nlink, UserLen, username, GroupLen, groupname,
            SizeLen, size, modTime, color+name+Reset, linked)
    }

    return result
}
