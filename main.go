// Count the lines of code in a directory
package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "v2.0.0"

// Flag to print loc by directory
var print_dir_flag = flag.Bool("d", false, "")

// Flag for directories to exclude
var exclude_dirs_flag = flag.String("ed", "", "")

// Flag for files to exclude
var exclude_files_flag = flag.String("ef", "", "")

// Flag for languages to exclude
var exclude_langs_flag = flag.String("el", "", "")

// Flag to print loc by file
var print_file_flag = flag.Bool("f", false, "")

// Flag to include dot directories
var include_dot_dir_flag = flag.Bool("id", false, "")

// Flag for max depth of subdirectories to search through
var max_search_depth = flag.Int("md", 1000, "")

// Flag for max number of file loc totals to print per directory
var max_print_files = flag.Int("mf", 1000000, "")

// Flag for max number of language loc totals to print per directory
var max_print_totals = flag.Int("ml", 1000, "")

// Flag to pick which column to sort by
var sort_column = flag.String("s", "loc", "")

// Flag to print version and exit
var version_flag = flag.Bool("v", false, "")

func main() {
	// Handle flags and arguments
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("flags/arguments not properly formatted")
		flag.Usage()
	}
	if *version_flag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

	// Set the directory that will be searched
	var dir_path string
	if len(args) == 1 {
		dir_path = args[0]
	} else {
		var err error
		dir_path, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting cwd:", err)
			return
		}
	}

	// Main functionality
	main_dir := newDirectory(dir_path, 0)
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
