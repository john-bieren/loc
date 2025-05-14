package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// Return path to parent of given entry
func parentDir(dir_path string) string {
	path_parts := splitPath(dir_path)
	parent_path_parts := path_parts[:len(path_parts)-1]
	parent_path := filepath.Join(parent_path_parts...)
	// the \ after the letter drive is not added by Join
	if runtime.GOOS == "windows" && parent_path[1] == ':' {
		parent_path = strings.Replace(parent_path, ":", ":\\", 1)
	}
	return parent_path
}

// Split filepath by slashes
func splitPath(path string) []string {
	path = strings.ReplaceAll(path, "\\", "/")
	return strings.Split(path, "/")
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
	path_parts := splitPath(full_path)
	rel_path_parts := path_parts[len(path_parts)-parents-1:]
	rel_path := filepath.Join(rel_path_parts...)
	return rel_path
}

// Remove duplicate values from a slice
func removeSliceDuplicates[T comparable](input_slice []T) []T {
	values := make(map[T]bool)
	unique_slice := []T{}
	for _, item := range input_slice {
		if _, exists := values[item]; !exists {
			values[item] = true
			unique_slice = append(unique_slice, item)
		}
	}
	return unique_slice
}

// Make sure no directory will be counted twice
func removeOverlappingDirs(dir_paths []string) []string {
	dir_paths = removeSliceDuplicates(dir_paths)
	var result []string
	for i, i_path := range dir_paths {
		keep_dir := true
		for j, j_path := range dir_paths {
			if i == j {
				continue
			}
			// if one path is contained within another
			if strings.Contains(i_path, j_path) {
				// drop the path unless -md dictates that it won't be searched otherwise
				i_split, j_split := splitPath(i_path), splitPath(j_path)
				distance := len(i_split) - len(j_split)
				if distance <= *max_search_depth {
					keep_dir = false
					break
				}
			}
		}
		if keep_dir {
			result = append(result, i_path)
		}
	}
	return result
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

// Make a slice of keys from a map[string]int sorted by value
func sortKeys(source_map map[string]int) []string {
	// extract keys from source map
	var keys []string
	for key := range source_map {
		keys = append(keys, key)
	}
	quickSort(source_map, keys, 0, len(keys)-1)
	return keys
}

// Sum the integer values of a map
func sumMapValues[k comparable](m map[k]int) int {
	var sum int
	for _, value := range m {
		sum += value
	}
	return sum
}

// Convert all paths to absolute paths
func toAbsPath(dir_paths []string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintln("Error getting cwd:", err))
	}

	for i, path := range dir_paths {
		if path == "." {
			dir_paths[i] = cwd
		} else if path == ".." {
			dir_paths[i] = parentDir(cwd)
		} else if !filepath.IsAbs(path) {
			if strings.HasPrefix(path, "-") {
				fmt.Printf("Warning: argument \"%s\" may be a misplaced flag, flags must come before arguments\n", path)
			}
			dir_paths[i] = filepath.Join(cwd, path)
		}
	}
	return dir_paths
}
