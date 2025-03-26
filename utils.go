package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

// Alphabetically sort string map keys
func alphaSortKeys[k any](m map[string]k) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
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

// Custom usage output for --help and relevant error messages
func usage() {
	fmt.Println("Count lines of code (loc) in a directory by language")
	fmt.Println("Usage: loc [options] [path]")
	fmt.Println("	options must come before path")
	fmt.Println("	path defaults to current working directory if no argument is given")
	fmt.Println("Options:")
	fmt.Println("	-d	Print loc totals by directory")
	fmt.Println("	-f	Print loc totals by file")
	fmt.Println("	-m int	Maximum depth of subdirectories to search")
	os.Exit(0)
}
