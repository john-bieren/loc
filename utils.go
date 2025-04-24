package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Convert integer to string with thousands separators
func addCommas(num int) string {
	str := strconv.Itoa(num)

	if len(str) <= 3 {
		return str
	}

	// create backwards string with commas
	var result []byte
	for count, i := 0, len(str)-1; i >= 0; count, i = count+1, i-1 {
		if count > 0 && count%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, str[i])
	}

	// reverse string
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

// Convert byte count into units
func formatByteCount(byte_count int) string {
	if byte_count <= 1_000 {
		return fmt.Sprintf("%d b", byte_count)
	} else if byte_count <= 1_000_000 {
		return fmt.Sprintf("%.1f kb", float64(byte_count)/1_000)
	} else if byte_count <= 1_000_000_000 {
		return fmt.Sprintf("%.1f mb", float64(byte_count)/1_000_000)
	} else {
		return fmt.Sprintf("%.1f gb", float64(byte_count)/1_000_000_000)
	}
}

// Quick sort implementation for sorting integer map values in descending order
func quickSort(source_map map[string]int, keys []string, low int, high int) {
	if low < high {
		// median of three to pick pivot value
		mid := low + (high-low)/2
		a, b, c := source_map[keys[high]], source_map[keys[mid]], source_map[keys[low]]
		var p int
		if (a > b) != (a > c) {
			p = high
		} else if (b > a) != (b > c) {
			p = mid
		} else {
			p = low
		}
		pivot := source_map[keys[p]]
		keys[p], keys[low] = keys[low], keys[p]
		i := high

		for j := high; j > low; j-- {
			if source_map[keys[j]] < pivot {
				keys[i], keys[j] = keys[j], keys[i]
				i--
			}
		}
		keys[i], keys[low] = keys[low], keys[i]

		quickSort(source_map, keys, low, i-1)
		quickSort(source_map, keys, i+1, high)
	}
}

// Convert full path to relative path from main_dir
func relPath(full_path string, parents int) string {
	full_path = strings.ReplaceAll(full_path, "\\", "/")
	path_parts := strings.Split(full_path, "/")
	rel_path_parts := path_parts[len(path_parts)-parents-1:]
	rel_path := filepath.Join(rel_path_parts...)
	return rel_path
}

// Sort a slice of files by loc or size
func sortFiles(slice []*file, sort_by string) []*file {
	sort.Slice(slice, func(i, j int) bool {
		if sort_by == "size" {
			return slice[i].bytes > slice[j].bytes
		} else {
			return slice[i].loc > slice[j].loc
		}
	})
	return slice
}

// Make a sorted slice of keys from a map[string]int using quick sort
func sortKeys(source_map map[string]int) []string {
	var keys []string
	for key := range source_map {
		keys = append(keys, key)
	}
	quickSort(source_map, keys, 0, len(keys)-1)
	return keys
}

// Sum the integer values of a map
func sumValues[k comparable](m map[k]int) int {
	var sum int
	for _, value := range m {
		sum += value
	}
	return sum
}

// Custom usage output for --help and relevant error messages
func usage() {
	fmt.Println("loc", version)
	fmt.Println("Count lines of code in directories and their subdirectories by language")
	fmt.Println("")
	fmt.Println("Usage: loc [options] [paths]")
	fmt.Println("         Options must come before paths")
	fmt.Println("         Paths are the directories you wish to search (cwd by default)")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("        -d        Print loc by directory")
	fmt.Println("        -ed str   Directories to exclude (use name or full path, i.e. \"src,lib,C:/Users/user/loc\")")
	fmt.Println("        -ef str   Files to exclude (use name or full path, i.e. \"index.js,utils.go,C:/Users/user/lib/main.py\")")
	fmt.Println("        -el str   Languages to exclude (i.e. \"HTML,Plain Text,JSON\")")
	fmt.Println("        -f        Print loc by file")
	fmt.Println("             -mf int   Maximum number of files to print per directory (default: 100,000)")
	fmt.Println("        -id       Include dot directories")
	fmt.Println("        -md int   Maximum depth of subdirectories to search (default: 1,000)")
	fmt.Println("        -ml int   Maximum number of language loc totals to print per directory (default: 1,000)")
	fmt.Println("        -p        Print loc as a percentage of overall total")
	fmt.Println("        -s str    Choose how to sort results [\"loc\", \"size\", \"files\"] (default: \"loc\")")
	fmt.Println("        -v        Print version and exit")
	os.Exit(0)
}
