package internals

import (
	"os"
	"syscall"
)

func calculateTotalBlocks(path string, includeHidden bool) int64 {
	var total int64

	entries, err := os.ReadDir(path)
	if err != nil {
		return total
	}
	for _, entry := range entries {
		if!includeHidden && isHidden(entry.Name()) {
            continue
        }
		name := entry.Name()
		names, err := os.Lstat(path + "/" + name)
		if err != nil {
			continue
		}
		blocksize := names.Sys().(*syscall.Stat_t)
		total += blocksize.Blocks / 2

	}
	return total
}

func isHidden(name string) bool {
	return name[0] == '.'
}
