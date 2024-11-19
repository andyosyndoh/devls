package internals

import (
	"fmt"
	"os"
)

func calculateTotalBlocks(path string, includeHidden bool) int64 {
	var total int64

	files, _ := os.ReadDir(path)
	for _, file := range files {
		if !includeHidden && file.Name()[0] == '.' {
			continue
		}
		info, err := file.Info()
		if err == nil {
			size := info.Size()
			blocks := (size + 511) / 512
			total += blocks
			fmt.Printf("Debug: File %s, size %d, blocks %d\n", file.Name(), size, blocks)
		}
	}

	if includeHidden {
		if info, err := os.Stat(path); err == nil {
			currentDirBlocks := (info.Size() + 511) / 512
			total += currentDirBlocks
			fmt.Printf("Debug: Current directory ., blocks %d\n", currentDirBlocks)
		}

		if parentInfo, err := os.Stat(dirName(path)); err == nil {
			parentDirBlocks := (parentInfo.Size() + 511) / 512
			total += parentDirBlocks
			fmt.Printf("Debug: Parent directory .., blocks %d\n", parentDirBlocks)
		}
	}

	// Convert to 1K blocks without rounding up
	totalKB := total / 2
	fmt.Printf("Debug: Total blocks before conversion: %d, after conversion: %d\n", total, totalKB)

	return totalKB
}
