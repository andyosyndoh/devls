package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

var colorMap = Colors()

const (
	Reset = "\033[0m" // Reset color
)

func main() {
	flags := map[string]bool{
		"a": false,
		"r": false,
		"t": false,
		"R": false,
		"l": false,
	}
	var filestore []string
	if len(os.Args) == 2 && os.Args[1] == "-" {
		// If so, just exit the program without doing anything
		return
	}
	// Parse command-line arguments for flags and target paths.
	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			for _, flag := range arg[1:] {
				switch flag {
				case 'a':
					flags["a"] = true
				case 'r':
					flags["r"] = true
				case 't':
					flags["t"] = true
				case 'R':
					flags["R"] = true
				case 'l':
					flags["l"] = true
				}
			}
		} else {
			filestore = append(filestore, arg)
		}
	}

	if len(filestore) == 0 { // Default to current directory if no path is provided.
		filestore = append(filestore, ".")
	}

	// Use the flags to determine the listing behavior
	if flags["l"] {
		longList(filestore, flags)
	} else {
		shortList(filestore, flags)
	}
}

func shortList(filestore []string, flags map[string]bool) {
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
		if !fileInfo.IsDir() {
			if flags["a"] || file[0] != '.' {
				validFiles = append(validFiles, file)
			}
		} else {
			if flags["R"] {
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
						directories = append(directories, entry)
					}
				}
				if flags["a"] {
					directories = append([]string{".", ".."}, directories...)
					for _, entry := range dirContents {
						if entry[0] == '.' && entry != "." && entry != ".." {
							directories = append(directories, entry)
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

func directoryList(dircontent []string, file string) []string {
	content, err := os.Open(file)
	if err != nil {
		return dircontent
	}

	names, err := content.Readdirnames(0)
	if err != nil {
		return dircontent
	}
	return names
}

func longList(files []string, flags map[string]bool) {
	// files = customAlphaSort(files)

	// Sort files based on flags before processing
	// if flags["t"] {
	// 	files = sortFilesByModTime(files)
	// } else{
	// 	files = customAlphaSort(files)
	// }

	// if flags["r"] {
	// 	for i := 0; i < len(files)/2; i++ {
	// 		j := len(files) - 1 - i
	// 		files[i], files[j] = files[j], files[i]
	// 	}
	// }

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

			// fmt.Println("After sorting:")
			// for _, entry := range entries {
			// 	fmt.Printf("  %s\n", entry.Name())
			// }

			totalBlocks := calculateTotalBlocks(file, flags["a"])
			fmt.Printf("Debug: Total blocks before division: %d\n", totalBlocks*2)
			fmt.Printf("total %d\n", totalBlocks)

			if flags["a"] {
				// Add . and .. at the beginning of the list
				dotEntries := []os.DirEntry{createDotEntry(".", file), createDotEntry("..", dirName(file))}
				entries = append(dotEntries, entries...)
			}

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
				fmt.Println()
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

type customDirEntry struct {
	name string
	path string
}

func createDotEntry(name string, path string) os.DirEntry {
	return &customDirEntry{name: name, path: path}
}

func (c *customDirEntry) Name() string {
	return c.name
}

func (c *customDirEntry) IsDir() bool {
	return true
}

func (c *customDirEntry) Type() os.FileMode {
	return os.ModeDir
}

func (c *customDirEntry) Info() (os.FileInfo, error) {
	return os.Stat(c.path)
}

// Sorting function to ensure lexicographical order.
// func customAlphaSort(files []string) []string {
// 	n := len(files)
// 	for i := 0; i < n-1; i++ {
// 		for j := 0; j < n-i-1; j++ {
// 			if strings.Compare(files[j], files[j+1]) > 0 {
// 				files[j], files[j+1] = files[j+1], files[j]
// 			}
// 		}
// 	}
// 	return files
// }

// A function that checks if the file exists and gets information about it.
func check(file string) (bool, fs.FileInfo, bool) {
	info, err := os.Lstat(file)
	if err != nil {
		return false, nil, false
	}
	isSymlink := info.Mode()&os.ModeSymlink != 0
	return true, info, isSymlink
}

// Function to determine if a file should be shown, taking '-a' into account.
func shouldShowFile(name string, showHidden bool) bool {
	if showHidden {
		return true
	}
	return !strings.HasPrefix(name, ".")
}

func sortEntries(entries []os.DirEntry, flags map[string]bool) {
	n := len(entries)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if compareEntries(entries[j+1], entries[j]) {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}

	if flags["r"] {
		reverseEntries(entries)
	}
	if flags["t"] {
		sortEntriesByTime(entries)
	}
}

func compareEntries(a, b os.DirEntry) bool {
	return compareStrings(a.Name(), b.Name())
}

func compareStrings(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	for i := 0; i < len(aRunes) && i < len(bRunes); i++ {
		aLower := unicode.ToLower(aRunes[i])
		bLower := unicode.ToLower(bRunes[i])
		if aLower != bLower {
			return aLower < bLower
		}
		if aRunes[i] != bRunes[i] {
			return aRunes[i] < bRunes[i]
		}
	}
	return len(aRunes) < len(bRunes)
}

func reverseEntries(entries []os.DirEntry) {
	for i := 0; i < len(entries)/2; i++ {
		j := len(entries) - 1 - i
		entries[i], entries[j] = entries[j], entries[i]
	}
}

// func isDigit(c rune) bool {
// 	return c >= '0' && c <= '9'
// }

// func isLetter(r rune) bool {
// 	// Check for ASCII letters
// 	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
// 		return true
// 	}
// 	return false
// }

// func isSpecialChar(c byte) bool {
// 	return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
// }

// func extractNumber(s string) (int, int) {
// 	var numStr string
// 	for i, c := range s {
// 		if !isDigit((c)) {
// 			if i == 0 {
// 				return 0, 0
// 			}
// 			num, _ := strconv.Atoi(numStr)
// 			return num, i
// 		}
// 		numStr += string(c)
// 	}
// 	num, _ := strconv.Atoi(numStr)
// 	return num, len(s)
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

func sortEntriesByTime(entries []os.DirEntry) {
	n := len(entries)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			infoI, _ := entries[j].Info()
			infoJ, _ := entries[j+1].Info()
			if infoI.ModTime().Before(infoJ.ModTime()) {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}
}

func getLongFormat(path string) string {
	//////fmt.Printf("Debug: Entering getLongFormat for path: %s\n", path)

	linkInfo, err := os.Lstat(path)
	if err != nil {
		// //fmt.Printf("Debug: Error with Lstat: %v\n", err)
		return ""
	}

	//////fmt.Printf("Debug: Lstat successful. File mode: %v\n", linkInfo.Mode())

	mode := linkInfo.Mode()
	nlink := linkInfo.Sys().(*syscall.Stat_t).Nlink
	uid := linkInfo.Sys().(*syscall.Stat_t).Uid
	gid := linkInfo.Sys().(*syscall.Stat_t).Gid
	size := linkInfo.Size()
	modTime := linkInfo.ModTime().Format("Jan  2 15:04")
	name := baseName(path)
	color := GetFileColor(linkInfo)

	// mt.Printf("Debug: Basic file info - Size: %d, ModTime: %s, Name: %s\n", size, modTime, name)

	username := strconv.FormatUint(uint64(uid), 10)
	groupname := strconv.FormatUint(uint64(gid), 10)

	if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		username = u.Username
	}
	if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupname = g.Name
	}

	//////fmt.Printf("Debug: Username: %s, Groupname: %s\n", username, groupname)

	modeStr := mode.String()

	if mode&os.ModeSymlink != 0 {
		// //fmt.Println("Debug: File is a symlink")
		modeStr = "l" + modeStr[1:] // Ensure the first character is 'l' for symlinks
		link, err := os.Readlink(path)
		if err == nil {
			name = fmt.Sprintf("%s -> %s", name, link)
			//////fmt.Printf("Debug: Symlink target: %s\n", link)
			// For symlinks, we'll keep the size of the link itself, not the target
		}
	}

	result := fmt.Sprintf("%s %2d %s %s %6d %s %s%s%s", modeStr, nlink, username, groupname, size, modTime, color, name, Reset)
	return result
}

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

func listRecursiveLong(path string, flags map[string]bool, indent string) {
	fmt.Printf("\n%s%s:\n", indent, path)
	totalBlocks := calculateTotalBlocks(path, flags["a"])
	fmt.Printf("%stotal %d\n", indent, totalBlocks)

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
		fmt.Printf("%s%s\n", indent, format)
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

func SortStringsAscending(slice []string) []string {
	n := len(slice)
	// Bubble sort algorithm
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Compare adjacent elements
			if slice[j] > slice[j+1] {
				// Swap if they are in the wrong order
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}

func SortStringsDescending(slice []string) []string {
	n := len(slice)
	// Bubble sort algorithm
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Compare adjacent elements
			if slice[j] < slice[j+1] {
				// Swap if they are in the wrong order
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}

func sortFilesByModTime(files []string) []string {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			infoI, errI := os.Stat(files[j])
			infoJ, errJ := os.Stat(files[j+1])
			if errI != nil || errJ != nil {
				if files[j] > files[j+1] {
					files[j], files[j+1] = files[j+1], files[j]
				}
			} else if infoI.ModTime().Before(infoJ.ModTime()) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
	return files
}

func print(files []string) {
	for _, value := range files {
		fmt.Println(value)
	}
}

func printShort(files []string) {
	var result string
	for i, value := range files {
		result += value
		if i < len(files) {
			result += "  "
		}
	}

	fmt.Println(result)
}

// Custom function to join path elements
func joinPath(elem ...string) string {
	return strings.Join(elem, "/")
}

// Custom function to get the base name of a path
func baseName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// Custom function to get the directory name of a path
func dirName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return "."
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

func GetFileColor(file os.FileInfo) string {
	if file.IsDir() {
		if color, ok := colorMap["di"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[34m" // Default to blue if not found
	}

	// Check for symbolic links (symlinks)
	if file.Mode()&os.ModeSymlink != 0 {
		if color, ok := colorMap["ln"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[0m" // Default to light cyan if not found
	}

	// Check for executable files
	if file.Mode().Perm()&0o111 != 0 {
		if color, ok := colorMap["ex"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[32m" // Default to green if executable color is not found
	}

	// Check for block devices
	if file.Mode()&os.ModeDevice != 0 && file.Mode()&os.ModeCharDevice == 0 {
		if color, ok := colorMap["bd"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[33m" // Default to yellow for block devices
	}

	// Check for character devices
	if file.Mode()&os.ModeCharDevice != 0 {
		if color, ok := colorMap["cd"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[33m" // Default to yellow for character devices
	}

	// Check for named pipes (e.g., FIFO files)
	if file.Mode()&os.ModeNamedPipe != 0 {
		if color, ok := colorMap["pi"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[31m" // Default to red if not found
	}

	// Fallback to reset if no specific color is found
	return Reset
}

func Colors() map[string]string {
	lsColors := os.Getenv("LS_COLORS")
	colorMap := make(map[string]string)

	if lsColors == "" {
		return colorMap // Return empty map if LS_COLORS is not set
	}

	pairs := strings.Split(lsColors, ":")
	for _, pair := range pairs {
		if strings.Contains(pair, "=") {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				colorMap[parts[0]] = parts[1]
			}
		}
	}

	return colorMap
}
