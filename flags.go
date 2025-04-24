package main

import "flag"

var (
	// Flag to print loc by directory
	print_dir_flag = flag.Bool("d", false, "")

	// Flag for directories to exclude
	exclude_dirs_flag = flag.String("ed", "", "")

	// Flag for files to exclude
	exclude_files_flag = flag.String("ef", "", "")

	// Flag for languages to exclude
	exclude_langs_flag = flag.String("el", "", "")

	// Flag to print loc by file
	print_file_flag = flag.Bool("f", false, "")

	// Flag to include dot directories
	include_dot_dir_flag = flag.Bool("id", false, "")

	// Flag for max depth of subdirectories to search through
	max_search_depth = flag.Int("md", 1_000, "")

	// Flag for max number of files to print per directory
	max_print_files = flag.Int("mf", 100_000, "")

	// Flag for max number of language loc totals to print per directory
	max_print_totals = flag.Int("ml", 1_000, "")

	// Flag to print loc as a percentage of overall total
	percentages_flag = flag.Bool("p", false, "")

	// Flag for which column to sort by
	sort_column = flag.String("s", "loc", "")

	// Flag to print version and exit
	version_flag = flag.Bool("v", false, "")
)
