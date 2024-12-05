package internals

import (
	"fmt"
	"os"
	"strings"
)

func LongList(files []string, flags map[string]bool) {
	files = sortFiles(files)
	for i, file := range files {
		exist, fileInfo, isSymlink := check(file)
		if !exist {
			fmt.Printf("ls: cannot access '%v': No such file or directory\n", file)
			continue
		}

		if isSymlink {
			format := getLongFormat(file, false)
			link, err := os.Readlink(file)
			_, fileInfo2, _ := check(link)
			if err == nil {
				linkColor := GetFileColor(fileInfo2.Mode(), link)
				color := GetFileColor(fileInfo.Mode(), file)
				fmt.Printf("%s %s%s%s -> %s%s%s\n", format, color, file, Reset, linkColor, link, Reset)
			} else {
				fmt.Println(format)
			}
			continue
		}

		if !fileInfo.IsDir() {
			if shouldShowFile(file, flags["a"]) {
				format := getLongFormat(file, false)

				color := GetFileColor(fileInfo.Mode(), file)
				fmt.Printf("%s %s\n", format, color+file+Reset)
				// fmt.Println(format)
			}
		} else {
			if len(files) > 1 || flags["R"] {
				fmt.Printf("%s:\n", file)
			}

			dirPath := file
			if file[len(file)-1] != '/' {
				dirPath += "/"
			}

			dirEntries, err := os.ReadDir(dirPath)
			if err != nil {
				fmt.Printf("Error reading directory %s: %v\n", file, err)
				continue
			}

			totalBlocks := calculateTotalBlocks(dirPath, flags["a"])
			fmt.Printf("total %d\n", totalBlocks)

			var entries []os.DirEntry

			// Add dot entries if -a flag is set
			if flags["a"] {
				entries = append(entries, createDotEntry(".", dirPath))
				entries = append(entries, createDotEntry("..", dirName(strings.TrimRight(dirPath, "/"))))
			}

			// Add regular entries to the slice
			for _, entry := range dirEntries {
				if shouldShowFile(entry.Name(), flags["a"]) {
					entries = append(entries, entry)
				}
			}

			// Sort entries
			sortEntries(entries, flags)

			// Process entries
			for _, entry := range entries {
				entryPath := joinPath(dirPath, entry.Name())
				format := getLongFormat(entryPath, entry.Name() == "." || entry.Name() == "..")

				// Add color based on file type
				color := GetFileColor(entry.Type(), entry.Name())

				// Handle symlinks
				if entry.Type()&os.ModeSymlink != 0 {
					link, err := os.Readlink(entryPath)
					_, fileInfo2, _ := check(file)
					if err == nil {
						linkColor := GetFileColor(fileInfo2.Mode(), link)
						fmt.Printf("%s %s%s%s -> %s%s%s\n", format, color, entry.Name(), Reset, linkColor, link, Reset)
					} else {
						fmt.Printf("%s %s%s%s\n", format, color, entry.Name(), Reset)
					}
				} else {
					// Regular files and directories
					fmt.Printf("%s %s%s%s\n", format, color, entry.Name(), Reset)
				}
			}

			// Handle recursive listing
			if flags["R"] {
				var subdirs []string
				for _, entry := range entries {
					if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
						subdir := joinPath(file, entry.Name())
						subdir = cleanPath(subdir)
						subdirs = append(subdirs, subdir)
					}
				}

				for j, subdir := range subdirs {
					fmt.Println()
					LongList([]string{subdir}, flags)
					if j < len(subdirs)-1 || i < len(files)-1 {
						// fmt.Println()
						// fmt.Println()
					}
				}
			}
		}

		if len(files) > 1 && i < len(files)-1 {
			exist, fileInfo, _ := check(file)
			if exist && fileInfo.IsDir() {
				fmt.Println()
			}
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
			fmt.Printf("ls: cannot access '%s': %v\n", file, err)
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
