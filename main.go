//go:generate go run ./generator
//go:generate gofmt -w languages.go

package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "v3.0.0 beta"

var total_loc, total_bytes, total_files float64

// main is loc's entry point
func main() {
	// overwrite default usage function so custom message prints on --help or argument errors
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	processFlags()

	var dir_paths []string
	if len(args) != 0 {
		dir_paths = args
	} else {
		cwd, err := os.Getwd()
		dir_paths = []string{cwd}
		if err != nil {
			panic(fmt.Sprintln("Error getting cwd:", err))
		}
	}

	dir_paths = toAbsPath(dir_paths)
	if len(dir_paths) > 1 {
		dir_paths = removeOverlappingDirs(dir_paths)
	}

	var main_dir *directory
	if len(dir_paths) == 1 {
		main_dir = newDirectory(dir_paths[0], 0)
	} else {
		// increment search depth since the "total" directory isn't real but counts as a parent
		*max_search_depth++

		// create a fake directory to show totals across multiple directory args
		main_dir = &directory{
			search_subdirs: true,
			loc_counts:     make(map[string]int),
			file_counts:    make(map[string]int),
			byte_counts:    make(map[string]int),
		}

		for _, path := range dir_paths {
			subdir := newDirectory(path, 1)
			main_dir.subdirectories = append(main_dir.subdirectories, subdir)
		}

		main_dir.countDirLoc()
	}

	if *percentages_flag {
		total_loc = float64(sumMapValues(main_dir.loc_counts))
		total_bytes = float64(sumMapValues(main_dir.byte_counts))
		total_files = float64(sumMapValues(main_dir.file_counts))
	}

	if len(main_dir.loc_counts) > 0 {
		if *print_dir_flag {
			main_dir.printTreeLoc()
		} else {
			main_dir.printDirLoc()
			if *print_file_flag {
				main_dir.printFileLoc()
			}
		}
	} else {
		fmt.Println("No code files found")
	}
}
