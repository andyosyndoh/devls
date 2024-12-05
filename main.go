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
