package internals

import (
	"fmt"
	"os"
	"strings"
)

func listRecursiveLong(path string, flags map[string]bool, indent string) {
	fmt.Printf("\n%s:\n", path)
	totalBlocks := calculateTotalBlocks(path, flags["a"])
	fmt.Printf("total %d\n", totalBlocks)

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", path, err)
		return
	}

	var filteredEntries []os.DirEntry
	if flags["a"] {
		// Include . and .. when -a flag is used
		filteredEntries = append(filteredEntries, createDotEntry(".", path), createDotEntry("..", dirName(path)))
		filteredEntries = append(filteredEntries, entries...)
	} else {
		for _, entry := range entries {
			if !strings.HasPrefix(entry.Name(), ".") {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	}

	sortEntries(filteredEntries, flags)

	for _, entry := range filteredEntries {
		entryPath := joinPath(path, entry.Name())
		if entry.Name() == "." {
			entryPath = path
		} else if entry.Name() == ".." {
			entryPath = dirName(path)
		}
		format := getLongFormat(entryPath)
		if entry.Name() == "." {
			format = strings.Replace(format, baseName(path), ".", 1)
		} else if entry.Name() == ".." {
			format = strings.Replace(format, baseName(dirName(path)), "..", 1)
		}
		fmt.Printf("%s\n", format)
	}

	var subdirs []string
	for _, entry := range filteredEntries {
		if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
			subdirs = append(subdirs, entry.Name())
		}
	}

	for _, subdir := range subdirs {
		fullPath := joinPath(path, subdir)
		listRecursiveLong(fullPath, flags, indent+"  ")
	}
}


func listRecursive(path string, flags map[string]bool, indent string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", path, err)
		return
	}

	fmt.Printf("%s%s:\n", indent, path)
	var entries []string
	for _, file := range files {
		if flags["a"] || file.Name()[0] != '.' {
			entries = append(entries, file.Name())
		}
	}

	if flags["t"] {
		entries = sortFilesByModTime(entries)
	} else if flags["r"] {
		entries = SortStringsDescending(entries)
	} else {
		entries = SortStringsAscending(entries)
	}

	
	printShort(entries)
	fmt.Println()

	for _, entry := range entries {
		fullPath := joinPath(path, entry)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if info.IsDir() {
			listRecursive(fullPath, flags, indent+"  ")
		}
	}
}