package internals

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)


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

func printShort(files []string, path string) {
	var result string
	for i, value := range files {
		var exist bool
		var fileInfo fs.FileInfo
		if path != "" {
			exist, fileInfo, _ = check(path + "/" + value)
		} else {
			exist, fileInfo, _ = check(value)
		}
		if !exist {
			continue
		}
		color := GetFileColor(fileInfo.Mode(), fileInfo.Name())
		result += color + value + Reset
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



// Custom function to get the directory name of a path
func dirName(path string) string {
    // Remove trailing slashes
	if path == "/dev/.." {
        return "/"
    }
    if strings.HasPrefix(path, "/dev/") {
        return "/dev"
    }

    // Rest of the original function...
    path = strings.TrimRight(path, "/")
    
    if path == "." {
        return ".."
    }

    if path == ".." {
        return "../.."
    }

    // Get absolute path if possible
    absPath, err := os.Getwd()
    if err == nil {
        fullPath := joinPath(absPath, path)
        // Split the path
        parts := strings.Split(fullPath, "/")
        if len(parts) > 1 {
            return strings.Join(parts[:len(parts)-1], "/")
        }
    }

    // Fallback to simple parent directory
    parts := strings.Split(path, "/")
    if len(parts) <= 1 {
        return ".."
    }
    return strings.Join(parts[:len(parts)-1], "/")
}