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
	children_loc   map[string]int
}

// Index children and subdirectories
func (d *directory) searchDir() {
	entries, err := os.ReadDir(d.full_path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "__") {
			continue
		}

		full_path := filepath.Join(d.full_path, entry.Name())
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
				child := newDirectory(full_path, d.parents+1)
				d.subdirectories = append(d.subdirectories, child)
			}
		} else {
			child := newFile(full_path, d.parents)
			if child.is_code {
				d.children = append(d.children, child)
			}
		}
	}
}

// Count lines of code for each language in all indexed children
func (d *directory) countDirLoc() {
	for _, child := range d.children {
		d.children_loc[child.file_type] += child.loc
	}

	if d.search_subs {
		for _, sub := range d.subdirectories {
			for file_type, loc := range sub.children_loc {
				d.children_loc[file_type] += loc
			}
		}
	}
}

// Print loc by file type for directory
func (d directory) printDirLoc() {
	if len(d.children_loc) == 0 {
		return
	}
	spaces := strings.Repeat("    ", d.parents)

	// Print directory name, if applicable
	if *print_dir_flag {
		fmt.Printf("%s%s/\n", spaces, d.name)
		spaces += " " // loc totals should have an extra space if dir names are printed
	}

	// Print loc total
	if len(d.children_loc) > 1 {
		fmt.Printf("%s%s loc\n", spaces, addCommas(sumValues(d.children_loc)))
	}

	// Print loc totals by file type
	keys := alphaSortKeys(d.children_loc)
	for _, file_type := range keys {
		fmt.Printf("%s%s %s loc\n", spaces, addCommas(d.children_loc[file_type]), file_type)
	}
}

// Print loc by file type for directory and subdirectories, include files if -f used
func (d directory) printTreeLoc() {
	d.printDirLoc()

	if *print_file_flag {
		spaces := strings.Repeat("    ", d.parents+1)
		for _, child := range d.children {
			fmt.Printf("%s%s loc - %s\n", spaces, addCommas(child.loc), child.name)
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
	for _, child := range d.children {
		fmt.Printf(" %s loc - %s\n", addCommas(child.loc), child.rel_path)
	}

	if d.search_subs {
		for _, sub := range d.subdirectories {
			sub.printFileLoc()
		}
	}
}

// Constructor for instances of directory struct
func newDirectory(path string, num_parents int) *directory {
	self := &directory{
		full_path:    path,
		name:         filepath.Base(path),
		parents:      num_parents,
		search_subs:  num_parents+1 <= *max_depth_flag,
		children_loc: make(map[string]int),
	}
	self.searchDir()
	self.countDirLoc()
	return self
}
