package internals

import (
	"fmt"
	"os"
)

func ShortList(files []string, flags map[string]bool) {
	showDirectoryNames := flags["R"] || len(files) > 1
	isFirstDirectory := true

	for _, file := range files {
		exist, fileInfo, isSymlink := check(file)
		if !exist {
			fmt.Printf("ls: cannot access '%v': No such file or directory\n", file)
			continue
		}
		if isSymlink {
			fmt.Println(file)
			continue
		}

		if !fileInfo.IsDir() {
			if shouldShowFile(file, flags["a"]) {
				fmt.Println(file)
			}
		} else {
			if showDirectoryNames {
				if !isFirstDirectory {
					fmt.Println()
				}
				fmt.Printf("%s:\n", file)
				isFirstDirectory = false
			}

			dirEntries, err := os.ReadDir(file)
			if err != nil {
				fmt.Printf("Error reading directory %s: %v\n", file, err)
				continue
			}

			var entries []os.DirEntry
			for _, entry := range dirEntries {
				if shouldShowFile(entry.Name(), flags["a"]) {
					entries = append(entries, entry)
				}
			}

			sortEntries(entries, flags)

			if flags["a"] {
				dotEntries := []os.DirEntry{createDotEntry(".", file), createDotEntry("..", dirName(file))}
				entries = append(dotEntries, entries...)
			}

			var fileNames []string
			for _, entry := range entries {
				fileNames = append(fileNames, entry.Name())
			}

			printShort(fileNames, file)

			if flags["R"] {
				for _, entry := range entries {
					if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
						subdir := joinPath(file, entry.Name())
						fmt.Println()
						ShortList([]string{subdir}, flags)
					}
				}
			}
		}
	}
}
