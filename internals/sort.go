package internals

import (
	"os"
	"unicode"
)

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