package internals

import "os"

type CustomDirEntry struct {
	name string
	path string
}

func createDotEntry(name string, path string) os.DirEntry {
	return &CustomDirEntry{name: name, path: path}
}

func (c *CustomDirEntry) Name() string {
	return c.name
}

func (c *CustomDirEntry) IsDir() bool {
	return true
}

func (c *CustomDirEntry) Type() os.FileMode {
	return os.ModeDir
}

func (c *CustomDirEntry) Info() (os.FileInfo, error) {
	return os.Stat(c.path)
}

var SizeLen, LinkLen, UserLen, GroupLen, MajorLen, MinorLen = 0, 0, 0, 0, 0, 0
