package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

var (
	// Flag to print loc by directory
	print_dir_flag = flag.Bool("d", false, "")

	// Flag for directories to exclude
	exclude_dirs_flag = flag.String("ed", "", "")
	// Parsed list of inputs for -ed flag
	exclude_dirs []string

	// Flag for files to exclude
	exclude_files_flag = flag.String("ef", "", "")
	// Parsed list of inputs for -ef flag
	exclude_files []string

	// Flag for languages to exclude
	exclude_langs_flag = flag.String("el", "", "")
	// Parsed list of inputs for -el flag
	exclude_langs []string

	// Flag to print loc by file
	print_file_flag = flag.Bool("f", false, "")

	// Flag to include dot directories
	include_dot_dir_flag = flag.Bool("id", false, "")

	// Flag for languages to include
	include_langs_flag = flag.String("il", "", "")
	// Parsed list of inputs for -il flag
	include_langs []string

	// Flag for max depth of subdirectories to search through
	max_search_depth = flag.Int("md", 1_000, "")

	// Flag for max number of files to print per directory
	max_print_files = flag.Int("mf", 100_000, "")

	// Flag for max number of language loc totals to print per directory
	max_print_totals = flag.Int("ml", 1_000, "")

	// Flag for max depth of subdirectories to print
	max_print_depth = flag.Int("mp", 1_000, "")

	// Flag to print loc as a percentage of overall total
	percentages_flag = flag.Bool("p", false, "")

	// Flag for which column to sort by
	sort_column = flag.String("s", "loc", "")

	// Flag to print version and exit
	version_flag = flag.Bool("v", false, "")
)

// Run exit flags, parse list flags, check for invalid inputs
func processFlags() {
	if *version_flag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

	if *exclude_dirs_flag != "" {
		exclude_dirs = strings.Split(*exclude_dirs_flag, ",")
	}
	if *exclude_files_flag != "" {
		exclude_files = strings.Split(*exclude_files_flag, ",")
	}
	if *exclude_langs_flag != "" {
		exclude_langs = strings.Split(*exclude_langs_flag, ",")
	}
	if *include_langs_flag != "" {
		include_langs = strings.Split(*include_langs_flag, ",")
	}

	if !slices.Contains([]string{"loc", "size", "files"}, *sort_column) {
		fmt.Printf("-s input \"%s\" is invalid, defaulting to \"loc\"\n", *sort_column)
	}
}
