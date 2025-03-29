package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type directory struct {
	full_path      string
	name           string
	parents        int
	search_subs    bool
	subdirectories []*directory
	children       []*file
	loc_counts     map[string]int
	file_counts    map[string]int
	byte_counts    map[string]int
}

// Index children and subdirectories
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
			// ignore errors from inaccessible dirs
			if os.IsNotExist(err) {
				fmt.Println("Cannot access directory entry:", err)
				continue
			}
			fmt.Println("Error checking directory entry:", err)
			return
		}

		if info.IsDir() {
			if d.search_subs {
				if !*include_dot_dir_flag && strings.HasPrefix(entry_name, ".") {
					continue
				}
				// check directory exclusions
				var skip_dir bool
				for excl := range strings.SplitSeq(*exclude_dirs_flag, ",") {
					if entry_name == excl {
						skip_dir = true
					} else if full_path == excl {
						skip_dir = true
					}
				}
				if skip_dir {
					continue
				}
				child := newDirectory(full_path, d.parents+1)
				d.subdirectories = append(d.subdirectories, child)
			}
		} else {
			// check file exclusions
			var skip_file bool
			for excl := range strings.SplitSeq(*exclude_files_flag, ",") {
				if entry_name == excl {
					skip_file = true
				} else if full_path == excl {
					skip_file = true
				}
			}
			if skip_file {
				continue
			}
			size := info.Size()
			child := newFile(full_path, d.parents, size)

			if child.is_code {
				// check language exclusions
				var skip_lang bool
				for excl := range strings.SplitSeq(*exclude_langs_flag, ",") {
					if child.file_type == excl {
						skip_lang = true
					}
				}
				if skip_lang {
					continue
				}
				d.children = append(d.children, child)
			}
		}
	}
}

// Count lines of code for each language in all indexed children
func (d *directory) countDirLoc() {
	for _, child := range d.children {
		d.loc_counts[child.file_type] += child.loc
		d.file_counts[child.file_type]++
		d.byte_counts[child.file_type] += child.bytes
	}

	if d.search_subs {
		for _, sub := range d.subdirectories {
			for file_type, loc := range sub.loc_counts {
				d.loc_counts[file_type] += loc
			}
			for file_type, n := range sub.file_counts {
				d.file_counts[file_type] += n
			}
			for file_type, b := range sub.byte_counts {
				d.byte_counts[file_type] += b
			}
		}
	}
}

// Print loc by file type for directory
func (d directory) printDirLoc() {
	if len(d.loc_counts) == 0 {
		return
	}
	indent := strings.Repeat("    ", d.parents)

	// Print directory name, if applicable
	if *print_dir_flag {
		fmt.Printf("%s%s/\n", indent, d.name)
		indent += " " // loc totals should have an extra space if dir names are printed
	}

	// Print column labels on first directory
	if d.parents == 0 {
		fmt.Printf("\033[1m%sLanguage: loc | bytes | files\033[0m\n", indent)
	}

	// Print loc total if multiple languages are present
	if len(d.loc_counts) > 1 {
		fmt.Printf(
			"%s%d langs: %s | %s | %s\n",
			indent, len(d.loc_counts),
			addCommas(sumValues(d.loc_counts)),
			addCommas(sumValues(d.byte_counts)),
			addCommas(sumValues(d.file_counts)),
		)
	}

	// Print loc totals by file type
	var keys []string
	switch *sort_column {
	case "bytes":
		keys = sortKeys(d.byte_counts)
	case "files":
		keys = sortKeys(d.file_counts)
	default:
		keys = sortKeys(d.loc_counts)
	}
	for i, file_type := range keys {
		// Print language total even if -ml=0 if there's only one language
		if i+1 > *max_print_totals && len(d.loc_counts) != 1 {
			break
		}
		fmt.Printf(
			"%s%s: %s | %s | %s\n",
			indent, file_type,
			addCommas(d.loc_counts[file_type]),
			addCommas(d.byte_counts[file_type]),
			addCommas(d.file_counts[file_type]),
		)
	}
}

// Print loc by file type for directory and subdirectories, include files if -f used
func (d directory) printTreeLoc() {
	d.printDirLoc()

	if *print_file_flag {
		indent := strings.Repeat("    ", d.parents+1)
		// Print column labels on first directory
		if d.parents == 0 && len(d.children) > 0 {
			fmt.Println("\033[1m    loc | bytes - file\033[0m")
		}
		for i, child := range sortFiles(d.children, *sort_column) {
			if i+1 > *max_print_files {
				break
			}
			fmt.Printf("%s%s | %s - %s\n", indent, addCommas(child.loc), addCommas(child.bytes), child.name)
		}
	}

	if d.search_subs {
		for _, sub := range d.subdirectories {
			sub.printTreeLoc()
		}
	}
}

// Print loc by file for all files counted
func (d directory) printFileLoc() {
	// collect all children to sort them
	var files []*file
	files = d.appendFiles(files)

	// Print column labels
	fmt.Println("\033[1m loc | bytes - file\033[0m")
	for i, file := range sortFiles(files, *sort_column) {
		if i+1 > *max_print_files {
			break
		}
		fmt.Printf(" %s | %s - %s\n", addCommas(file.loc), addCommas(file.bytes), file.rel_path)
	}
}

// Append files from d.children to input slice
func (d directory) appendFiles(files []*file) []*file {
	files = append(files, d.children...)
	if d.search_subs {
		for _, sub := range d.subdirectories {
			files = sub.appendFiles(files)
		}
	}
	return files
}

// Constructor for instances of directory struct
func newDirectory(path string, num_parents int) *directory {
	self := &directory{
		full_path:   path,
		name:        filepath.Base(path),
		parents:     num_parents,
		search_subs: num_parents+1 <= *max_search_depth,
		loc_counts:  make(map[string]int),
		file_counts: make(map[string]int),
		byte_counts: make(map[string]int),
	}
	self.searchDir()
	self.countDirLoc()
	return self
}
