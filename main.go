// Count the lines of code in a directory
package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "v2.2.0 beta"

var total_loc, total_bytes, total_files float64

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if *version_flag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

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
	dir_paths = convertSpecialPaths(dir_paths)

	var main_dir *directory
	if len(dir_paths) == 1 {
		main_dir = newDirectory(dir_paths[0], 0)
	} else {
		// increment search depth since the "total" directory isn't real but counts as a parent
		*max_search_depth++

		// create a fake directory to show totals across multiple directory args
		main_dir = &directory{
			name:        "total",
			search_subs: true,
			loc_counts:  make(map[string]int),
			file_counts: make(map[string]int),
			byte_counts: make(map[string]int),
		}

		for _, path := range dir_paths {
			child := newDirectory(path, 1)
			main_dir.subdirectories = append(main_dir.subdirectories, child)
		}

		main_dir.countDirLoc()
	}

	if *percentages_flag {
		total_loc = float64(sumValues(main_dir.loc_counts))
		total_bytes = float64(sumValues(main_dir.byte_counts))
		total_files = float64(sumValues(main_dir.file_counts))
	}

	if len(main_dir.children)+len(main_dir.subdirectories) > 0 {
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
