package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// tree_file_headers_printed tracks whether file column headers have been printed by printTreeLoc.
var tree_file_headers_printed bool

type directory struct {
	full_path      string
	name           string
	parents        int
	search_subdirs bool
	subdirectories []*directory
	files          []*file
	loc_counts     map[string]int
	file_counts    map[string]int
	byte_counts    map[string]int
}

// searchDir indexes the directory's files and subdirectories.
func (d *directory) searchDir() {
	entries, err := os.ReadDir(d.full_path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		entry_name := entry.Name()
		full_path := filepath.Join(d.full_path, entry_name)
		info, err := os.Stat(full_path)
		if err != nil {
			// specify errors from inaccessible entries, a common case
			if os.IsNotExist(err) {
				fmt.Println("Cannot access directory entry:", err)
			} else {
				fmt.Println("Error checking directory entry:", err)
			}
			continue
		}

		if info.IsDir() {
			if d.search_subdirs {
				if !*include_dot_dir_flag && strings.HasPrefix(entry_name, ".") {
					continue
				}

				var skip_dir bool
				for _, excl := range exclude_dirs {
					if entry_name == excl || full_path == excl {
						skip_dir = true
						break
					}
				}
				if skip_dir {
					continue
				}

				subdir := newDirectory(full_path, d.parents+1)
				d.subdirectories = append(d.subdirectories, subdir)
			}
		} else {
			var skip_file bool
			for _, excl := range exclude_files {
				if entry_name == excl || full_path == excl {
					skip_file = true
					break
				}
			}
			if skip_file {
				continue
			}

			size := info.Size()
			file := newFile(full_path, d.parents, size)
			if file.is_code {
				d.files = append(d.files, file)
			}
		}
	}
}

// countDirLoc counts the lines of code for each language in all indexed files.
func (d *directory) countDirLoc() {
	for _, file := range d.files {
		d.loc_counts[file.file_type] += file.loc
		d.file_counts[file.file_type]++
		d.byte_counts[file.file_type] += file.bytes
	}

	if d.search_subdirs {
		for _, subdir := range d.subdirectories {
			for file_type, loc := range subdir.loc_counts {
				d.loc_counts[file_type] += loc
			}
			for file_type, n := range subdir.file_counts {
				d.file_counts[file_type] += n
			}
			for file_type, b := range subdir.byte_counts {
				d.byte_counts[file_type] += b
			}
		}
	}
}

// printDirLoc prints loc by file type for the directory.
func (d directory) printDirLoc() {
	if len(d.loc_counts) == 0 {
		return
	}
	indent := strings.Repeat("    ", d.parents)

	// print directory name, if applicable
	if *print_dir_flag && d.parents > 0 {
		fmt.Printf("%s%s/\n", indent, d.name)
		indent += " " // loc totals should have an extra space if dir names are printed
	}

	// print column labels on first directory
	if d.parents == 0 {
		fmt.Printf("\033[1m%sLanguage: loc | size | files\033[0m\n", indent)
	}

	// print loc total if multiple languages are present
	if len(d.loc_counts) > 1 {
		if *percentages_flag && d.parents > 0 {
			fmt.Printf(
				"%s%d langs: %.1f%% | %.1f%% | %.1f%%\n",
				indent, len(d.loc_counts),
				float64(sumMapValues(d.loc_counts))/total_loc*100,
				float64(sumMapValues(d.byte_counts))/total_bytes*100,
				float64(sumMapValues(d.file_counts))/total_files*100,
			)
		} else {
			fmt.Printf(
				"%s%d langs: %s | %s | %s\n",
				indent, len(d.loc_counts),
				addCommas(sumMapValues(d.loc_counts)),
				formatByteCount(sumMapValues(d.byte_counts)),
				addCommas(sumMapValues(d.file_counts)),
			)
		}
	}

	// print loc totals by file type
	var keys []string
	switch *sort_column {
	case "size":
		keys = sortKeys(d.byte_counts)
	case "files":
		keys = sortKeys(d.file_counts)
	default:
		keys = sortKeys(d.loc_counts)
	}
	for i, file_type := range keys {
		// print language total even if -ml=0 if there's only one language
		if i+1 > *max_print_totals && len(d.loc_counts) > 1 {
			break
		}
		if *percentages_flag && !(len(d.loc_counts) == 1 && d.parents == 0) {
			fmt.Printf(
				"%s%s: %.1f%% | %.1f%% | %.1f%%\n",
				indent, file_type,
				float64(d.loc_counts[file_type])/total_loc*100,
				float64(d.byte_counts[file_type])/total_bytes*100,
				float64(d.file_counts[file_type])/total_files*100,
			)
		} else {
			fmt.Printf(
				"%s%s: %s | %s | %s\n",
				indent, file_type,
				addCommas(d.loc_counts[file_type]),
				formatByteCount(d.byte_counts[file_type]),
				addCommas(d.file_counts[file_type]),
			)
		}
	}
}

// printTreeLoc prints loc by file type for the directory and its subdirectories, includes files if -f used.
func (d directory) printTreeLoc() {
	d.printDirLoc()

	if *print_file_flag {
		indent := strings.Repeat("    ", d.parents+1)
		if !tree_file_headers_printed && len(d.files) > 0 {
			fmt.Printf("\033[1m%sloc | size - file\033[0m\n", indent)
			tree_file_headers_printed = true
		}

		for i, file := range sortFiles(d.files, *sort_column) {
			if i+1 > *max_print_files {
				break
			}
			if *percentages_flag {
				fmt.Printf(
					"%s%.1f%% | %.1f%% - %s\n",
					indent,
					float64(file.loc)/total_loc*100,
					float64(file.bytes)/total_bytes*100,
					file.name,
				)
			} else {
				fmt.Printf(
					"%s%s | %s - %s\n",
					indent,
					addCommas(file.loc),
					formatByteCount(file.bytes),
					file.name,
				)
			}
		}
	}

	if d.search_subdirs && *max_print_depth >= d.parents+1 {
		// sort the subdirectories by the selected sort column
		sort.Slice(d.subdirectories, func(i, j int) bool {
			switch *sort_column {
			case "size":
				return sumMapValues(d.subdirectories[i].byte_counts) > sumMapValues(d.subdirectories[j].byte_counts)
			case "files":
				return sumMapValues(d.subdirectories[i].file_counts) > sumMapValues(d.subdirectories[j].file_counts)
			default:
				return sumMapValues(d.subdirectories[i].loc_counts) > sumMapValues(d.subdirectories[j].loc_counts)
			}
		})

		for _, subdir := range d.subdirectories {
			subdir.printTreeLoc()
		}
	}
}

// printFileLoc prints loc by file for all files counted.
func (d directory) printFileLoc() {
	// files is a slice of the directory's files used to sort them.
	var files []*file
	files = d.appendFiles(files)

	fmt.Println("\033[1m loc | size - file\033[0m")
	for i, file := range sortFiles(files, *sort_column) {
		if i+1 > *max_print_files {
			break
		}
		if *percentages_flag {
			fmt.Printf(
				" %.1f%% | %.1f%% - %s\n",
				float64(file.loc)/total_loc*100,
				float64(file.bytes)/total_bytes*100,
				file.rel_path,
			)
		} else {
			fmt.Printf(
				" %s | %s - %s\n",
				addCommas(file.loc),
				formatByteCount(file.bytes),
				file.rel_path,
			)
		}
	}
}

// appendFiles appends files from d.files to input slice.
func (d directory) appendFiles(files []*file) []*file {
	files = append(files, d.files...)
	if d.search_subdirs {
		for _, subdir := range d.subdirectories {
			files = subdir.appendFiles(files)
		}
	}
	return files
}

// newDirectory is the constructor for instances of the directory struct.
func newDirectory(path string, num_parents int) *directory {
	self := &directory{
		full_path:      path,
		name:           filepath.Base(path),
		parents:        num_parents,
		search_subdirs: num_parents+1 <= *max_search_depth,
		loc_counts:     make(map[string]int),
		file_counts:    make(map[string]int),
		byte_counts:    make(map[string]int),
	}
	self.searchDir()
	self.countDirLoc()
	return self
}
