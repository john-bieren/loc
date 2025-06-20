package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// addCommas converts an integer to a string with thousands separators.
func addCommas(num int) string {
	str := strconv.Itoa(num)
	if len(str) <= 3 {
		return str
	}

	// result is a reversed copy of str with commas.
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

// formatByteCount converts a raw byte count into a string formatted in the relveant units.
func formatByteCount(byteCount int) string {
	if byteCount <= 1_000 {
		return fmt.Sprintf("%d b", byteCount)
	} else if byteCount <= 1_000_000 {
		return fmt.Sprintf("%.1f kb", float64(byteCount)/1_000)
	} else if byteCount <= 1_000_000_000 {
		return fmt.Sprintf("%.1f mb", float64(byteCount)/1_000_000)
	} else {
		return fmt.Sprintf("%.1f gb", float64(byteCount)/1_000_000_000)
	}
}

// parentDir returns the path to the parent of the given entry.
func parentDir(dirPath string) string {
	pathParts := splitPath(dirPath)
	parentPathParts := pathParts[:len(pathParts)-1]
	parentPath := filepath.Join(parentPathParts...)

	// the \ after the drive letter is not added by Join
	if runtime.GOOS == "windows" && parentPath[1] == ':' {
		parentPath = strings.Replace(parentPath, ":", ":\\", 1)
	}
	return parentPath
}

// quickSort is a quick sort implementation for sorting map keys by their integer values in descending order.
func quickSort(sourceMap map[string]int, keys []string, low, high int) {
	if low < high {
		// median of three to pick pivot value
		mid := low + (high-low)/2
		a, b, c := sourceMap[keys[high]], sourceMap[keys[mid]], sourceMap[keys[low]]
		var p int
		if (a > b) != (a > c) {
			p = high
		} else if (b > a) != (b > c) {
			p = mid
		} else {
			p = low
		}
		pivot := sourceMap[keys[p]]
		keys[p], keys[low] = keys[low], keys[p]
		i := high

		for j := high; j > low; j-- {
			if sourceMap[keys[j]] < pivot {
				keys[i], keys[j] = keys[j], keys[i]
				i--
			}
		}
		keys[i], keys[low] = keys[low], keys[i]

		quickSort(sourceMap, keys, low, i-1)
		quickSort(sourceMap, keys, i+1, high)
	}
}

// removeSliceDuplicates removes duplicate values from a slice.
func removeSliceDuplicates[T comparable](inputSlice []T) []T {
	values := make(map[T]bool)
	uniqueSlice := []T{}
	for _, item := range inputSlice {
		if _, exists := values[item]; !exists {
			values[item] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}
	return uniqueSlice
}

// removeOverlappingDirs removes from a slice paths that are contained within other paths in the slice.
func removeOverlappingDirs(dirPaths []string) []string {
	dirPaths = removeSliceDuplicates(dirPaths)
	var result []string
	for i, iPath := range dirPaths {
		keepDir := true
		for j, jPath := range dirPaths {
			if i == j {
				continue
			}
			// if one path is contained within another
			if strings.Contains(iPath, jPath) {
				// drop the path unless -md dictates that it won't be searched otherwise
				iSplit, jSplit := splitPath(iPath), splitPath(jPath)
				distance := len(iSplit) - len(jSplit)
				if distance <= *maxSearchDepth {
					keepDir = false
					break
				}
			}
		}
		if keepDir {
			result = append(result, iPath)
		}
	}
	return result
}

// sortFiles sorts a slice of files by loc or size.
func sortFiles(slice []*file, sortBy string) []*file {
	sort.Slice(slice, func(i, j int) bool {
		if sortBy == "size" {
			return slice[i].bytes > slice[j].bytes
		} else {
			return slice[i].loc > slice[j].loc
		}
	})
	return slice
}

// sortKeys makes a slice of keys from a map[string]int, sorted by value.
func sortKeys(sourceMap map[string]int) []string {
	// keys is a copy of sourceMap's keys.
	var keys []string
	for key := range sourceMap {
		keys = append(keys, key)
	}
	quickSort(sourceMap, keys, 0, len(keys)-1)
	return keys
}

// splitPath splits a filepath by slashes.
func splitPath(path string) []string {
	path = strings.ReplaceAll(path, "\\", "/")
	return strings.Split(path, "/")
}

// sumMapValues sums the integer values of a map.
func sumMapValues[k comparable](m map[k]int) int {
	var sum int
	for _, value := range m {
		sum += value
	}
	return sum
}

// toAbsPath converts a path into an absolute path.
func toAbsPath(path string) string {
	if path == "." {
		path = cwd
	} else if path == ".." {
		path = parentDir(cwd)
	} else if !filepath.IsAbs(path) {
		if strings.HasPrefix(path, "-") {
			fmt.Printf("Warning: argument \"%s\" may be a misplaced flag, flags must come before arguments\n", path)
		}
		path = filepath.Join(cwd, path)
	}
	return path
}

// warn prints the message for a non-critical error if -q is not used.
func warn(message string, err error) {
	if !*suppressWarningsFlag {
		fmt.Println(message, err)
	}
}
