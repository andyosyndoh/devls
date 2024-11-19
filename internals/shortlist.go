package internals

import "fmt"

func ShortList(filestore []string, flags map[string]bool) {
	var message []string
	var validFiles []string
	var directories []string
	var errorMessage string
	var files []string
	if flags["t"] {
		files = sortFilesByModTime(filestore)
	} else if flags["r"] {
		files = SortStringsDescending(filestore)
	} else {
		files = SortStringsAscending(filestore)
	}
	for _, file := range files {
		exist, fileInfo, _ := check(file)
		if !exist {
			errorMessage = fmt.Sprintf("ls: cannot access '%v': No such file or directory", file)
			message = append(message, errorMessage)
			continue
		}
		color := GetFileColor(fileInfo.Mode(), fileInfo.Name())
		if !fileInfo.IsDir() {
			if flags["a"] || file[0] != '.' {
				validFiles = append(validFiles, color+file+Reset)
			}
		} else {
			if flags["R"] {
				fmt.Println("here")
				listRecursive(file, flags, "")
			} else {
				dirContents := directoryList([]string{}, file)
				if flags["t"] {
					dirContents = sortFilesByModTime(dirContents)
				} else if flags["r"] {
					dirContents = SortStringsDescending(dirContents)
				} else {
					dirContents = SortStringsAscending(dirContents)
				}
				for _, entry := range dirContents {
					if flags["a"] || entry[0] != '.' {
						ok, info, _ := check(entry)
						color := ""
						if ok {
							color = GetFileColor(info.Mode(), info.Name())
						}
						directories = append(directories, color+entry+Reset)
					}
				}
				if flags["a"] {
					directories = append([]string{".", ".."}, directories...)
					for _, entry := range dirContents {
						if entry[0] == '.' && entry != "." && entry != ".." {
							ok, info, _ := check(entry)
							color := ""
							if ok {
								color = GetFileColor(info.Mode(), info.Name())
							}
							directories = append(directories, color+entry+Reset)
						}
					}
				}
			}
		}
	}
	print(message)
	if len(validFiles) > 0 {
		printShort(validFiles)
	}
	if len(validFiles) > 0 && len(directories) > 0 && !flags["R"] {
		fmt.Println()
	}
	if !flags["R"] {
		printShort(directories)
	}
}
