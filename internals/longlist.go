package internals

import (
	"fmt"
	"os"
	"strings"
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
			if len(files) > 1 || flags["R"] {
				fmt.Printf("%s:\n", file)
			}
			if file[len(file)-1] != '/' {
				file += "/"
			}
			dirEntries, err := os.ReadDir(file)
			if err != nil {
				fmt.Printf("Error reading directory %s: %v\n", file, err)
				continue
			}

			var entries []os.DirEntry
			if flags["a"] {
				entries = append(entries, createDotEntry(".", file), createDotEntry("..", dirName(file)))
			}
			for _, entry := range dirEntries {
				if shouldShowFile(entry.Name(), flags["a"]) {
					entries = append(entries, entry)
				}
			}

			sortEntries(entries, flags)

			totalBlocks := calculateTotalBlocks(file, flags["a"])
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
				var subdirs []string
				for _, entry := range entries {
					if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
						subdir := joinPath(file, entry.Name())
						subdir = cleanPath(subdir) // Remove double slashes
						subdirs = append(subdirs, subdir)
					}
				}

				for j, subdir := range subdirs {
					fmt.Println() // Add a newline before each subdirectory listing
					LongList([]string{subdir}, flags)
					if j < len(subdirs)-1 || i < len(files)-1 {
						fmt.Println() // Add a newline after each subdirectory listing, except for the last one
					}
				}
			}
		}

		if len(files) > 1 && i < len(files)-1 {
			fmt.Println()
		}
	}
}

// Helper function to clean the path (remove double slashes)
func cleanPath(path string) string {
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	return path
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
