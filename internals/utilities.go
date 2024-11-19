package internals

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

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
