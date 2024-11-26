package internals

import (
	"fmt"
	"os"
)

func LongList(files []string, flags map[string]bool) {
	files = sortFiles(files)
	for i, file := range files {
		n := file
		exist, fileInfo, isSymlink := check(file)
		if !exist {
			fmt.Printf("ls: cannot access '%v': No such file or directory\n", file)
			continue
		}
		_ = calculateTotalBlocks(".", flags["a"])
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
				fmt.Printf("%v:\n", file)
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
				if n != "." && entry.Name() == "." || n != "." && entry.Name() == ".." {
					format = Name(format, entry.Name())
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

		if i != len(files)-1 {
			exist, fileInfo, _ := check(files[i+1])
			if exist && fileInfo.IsDir() {
				fmt.Println()
			}
		}
	}
}

func Name(format string, name string) string {
	color, reset := ("\033[" + lsColors["di"] + "m"), Reset
	format = format[:len(format)-1] + color + name + reset
	return format
}

func sortFiles(files []string) []string {
	// Separate files and directories
	var fileList []string
	var dirList []string

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", file, err)
			continue
		}
		if info.IsDir() {
			dirList = append(dirList, file)
		} else {
			fileList = append(fileList, file)
		}
	}

	// Sort both lists alphabetically
	SortStringsAscending(fileList)
	SortStringsAscending(dirList)

	// Combine files and directories, with files first
	return append(fileList, dirList...)
}
