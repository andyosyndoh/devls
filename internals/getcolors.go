package internals

import (
	"os"
	"strings"
)

var lsColors = Colors()

const (
	Reset = "\033[0m" // Reset color
)

// func GetFileColor(file os.FileInfo) string {
// 	if file.IsDir() {
// 		if color, ok := colorMap["di"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[34m" // Default to blue if not found
// 	}

// 	// Check for symbolic links (symlinks)
// 	if file.Mode()&os.ModeSymlink != 0 {
// 		if color, ok := colorMap["ln"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[0m" // Default to light cyan if not found
// 	}

// 	// Check for executable files
// 	if file.Mode().Perm()&0o111 != 0 {
// 		if color, ok := colorMap["ex"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[32m" // Default to green if executable color is not found
// 	}

// 	// Check for block devices
// 	if file.Mode()&os.ModeDevice != 0 && file.Mode()&os.ModeCharDevice == 0 {
// 		if color, ok := colorMap["bd"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[33m" // Default to yellow for block devices
// 	}

// 	// Check for character devices
// 	if file.Mode()&os.ModeCharDevice != 0 {
// 		if color, ok := colorMap["cd"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[33m" // Default to yellow for character devices
// 	}

// 	// Check for named pipes (e.g., FIFO files)
// 	if file.Mode()&os.ModeNamedPipe != 0 {
// 		if color, ok := colorMap["pi"]; ok {
// 			return "\033[" + color + "m"
// 		}
// 		return "\033[31m" // Default to red if not found
// 	}

// 	// Fallback to reset if no specific color is found
// 	return Reset
// }

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

func GetFileColor(mode os.FileMode, fileName string) string {
	switch {
	case mode&os.ModeDir != 0:
		if mode&0o002 != 0 && mode&0o010 != 0 {
			return "\033[" + lsColors["tw"] + "m" // Directory, writable by others, with sticky bit
		}
		if mode&0o002 != 0 {
			return "\033[" + lsColors["ow"] + "m" // Directory, writable by others
		}
		return "\033[" + lsColors["di"] + "m" // Directory
	case mode&os.ModeSymlink != 0:
		return "\033[" + lsColors["ln"] + "m" // Symlink
	case mode&os.ModeNamedPipe != 0:
		return "\033[" + lsColors["pi"] + "m" // Named pipe
	case mode&os.ModeSocket != 0:
		return "\033[" + lsColors["so"] + "m" // Socket
	case mode&os.ModeDevice != 0:
		return "\033[" + lsColors["bd"] + "m" // Block device
	case mode&os.ModeCharDevice != 0:
		return "\033[" + lsColors["cd"] + "m" // Character device
	case mode&os.ModeSetuid != 0:
		return "\033[" + lsColors["su"] + "m" // Setuid
	case mode&os.ModeSetgid != 0:
		return "\033[" + lsColors["sg"] + "m" // Setgid
	case mode&0o111 != 0:
		return "\033[" + lsColors["ex"] + "m" // Executable
	default:
		return getColorByExtension(strings.ToLower(getFileExtension(fileName)))
	}
}
func getColorByExtension(ext string) string {
	if color, ok := lsColors["*."+ext]; ok {
		return "\033[" + color + "m" 
	}
	return "\033[" + lsColors["rs"] + "m" // Default color
}

// getFileExtension extracts and returns the file extension from the given file name.
// If no extension is found, it returns an empty string.
func getFileExtension(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}