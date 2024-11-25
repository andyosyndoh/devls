package internals

import (
	"fmt"
	"os"
	"strings"
)

func LongList(files []string, flags map[string]bool) {
	for _, file := range files {
		exist, fileInfo, isSymlink := check(file)
		if !exist {
			fmt.Printf("ls: cannot access '%v': No such file or directory\n", file)
			continue
		}
		if isSymlink {
			format := getLongFormat(file)
			fmt.Println(format)
			continue
		}

		if !fileInfo.IsDir() {
			if shouldShowFile(file, flags["a"]) {
				format := getLongFormat(file)
				fmt.Println(format)
			}
		} else {
			if len(files) > 1 {
				fmt.Printf("%s:\n", file)
			}
			dirEntries, err := os.ReadDir(file)
			if err != nil {
				fmt.Printf("Error reading directory %s: %v\n", file, err)
				continue
			}

			var entries []os.DirEntry
            if flags["a"] {
                // Add . and .. at the beginning of the list
                entries = append(entries, createDotEntry(".", file), createDotEntry("..", dirName(file)))
            }
            for _, entry := range dirEntries {
                if shouldShowFile(entry.Name(), flags["a"]) {
                    entries = append(entries, entry)
                }
            }

			// fmt.Println("Before sorting:")
			// for _, entry := range entries {
			// 	fmt.Printf("  %s\n", entry.Name())
			// }

			sortEntries(entries, flags)

			if flags["R"] {
				fmt.Println(".:")
			}
			totalBlocks := calculateTotalBlocks(file, flags["a"])
			// fmt.Printf("Debug: Total blocks before division: %d\n", totalBlocks*2)
			fmt.Printf("total %d\n", totalBlocks)
			
			for _, entry := range entries {
				entryPath := joinPath(file, entry.Name())
				if entry.Name() == "." {
					entryPath = file
				} else if entry.Name() == ".." {
					entryPath = dirName(file)
				}
				format := getLongFormat(entryPath)
				if entry.Name() == "." {
					format = strings.Replace(format, baseName(file), ".", 1)
				} else if entry.Name() == ".." {
					format = strings.Replace(format, baseName(dirName(file)), "..", 1)
				}
				fmt.Println(format)
			}

			if flags["R"] {
				for _, entry := range entries {
					if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
						subdir := joinPath(file, entry.Name())
						listRecursiveLong(subdir, flags, "  ")
					}
				}
			}
		}

		if len(files) > 1 {
			fmt.Println()
		}
	}
}