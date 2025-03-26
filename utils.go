package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Sum the integer values of a map
func sumValues[k comparable](m map[k]int) int {
	var sum int
	for _, value := range m {
		sum += value
	}
	return sum
}

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

// Convert full path to relative path from main_dir
func relPath(full_path string, parents int) string {
	full_path = strings.ReplaceAll(full_path, "\\", "/")
	path_parts := strings.Split(full_path, "/")
	rel_path_parts := path_parts[len(path_parts)-parents-1:]
	rel_path := filepath.Join(rel_path_parts...)
	return rel_path
}

func sortFiles(slice []*file) []*file {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].loc > slice[j].loc
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

// Quick sort implementation for sorting integer map values in descending order
func quickSort(source_map map[string]int, keys []string, low int, high int) {
	if low < high {
		pivot := source_map[keys[low]]
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

// Custom usage output for --help and relevant error messages
func usage() {
	fmt.Println("loc", version)
	fmt.Println("Count lines of code (loc) in a directory by language")
	fmt.Println("")
	fmt.Println("Usage: loc [options] [path]")
	fmt.Println("	Options must come before path")
	fmt.Println("	Path defaults to current working directory if no argument is given")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("	-d	Print loc totals by directory")
	fmt.Println("	-f	Print loc totals by file")
	fmt.Println("	-m int	Maximum depth of subdirectories to search")
	fmt.Println("	-v	Print version and exit")
	os.Exit(0)
}
