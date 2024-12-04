package main

import (
	"ls/internals"
	"os"
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
		if arg[0] == '-' && len(arg) > 1 {
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
		internals.LongList(filestore, flags)
	} else {
		internals.ShortList(filestore, flags)
	}
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
