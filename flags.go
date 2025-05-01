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

	// Flag for max number of files to print per directory
	max_print_files = flag.Int("mf", 100_000, "")

	// Flag for max number of language loc totals to print per directory
	max_print_totals = flag.Int("ml", 1_000, "")

	// Flag to print loc as a percentage of overall total
	percentages_flag = flag.Bool("p", false, "")

	// Flag for max depth of subdirectories to print
	max_print_depth = flag.Int("pd", 1_000, "")

	// Flag for which column to sort by
	sort_column = flag.String("s", "loc", "")

	// Flag for max depth of subdirectories to search through
	max_search_depth = flag.Int("sd", 1_000, "")

	// Flag to print license information and exit
	license_flag = flag.Bool("license", false, "")

	// Flag to print version and exit
	version_flag = flag.Bool("version", false, "")

	// Line-by-line text for usage function output
	usage_message = []string{
		fmt.Sprintf("loc %s", version),
		"Count lines of code in directories and their subdirectories by language",
		"         Note: multi-line comments and non-comment docstrings are counted as lines of code",
		"",
		"Usage: loc [options] [dirs]",
		"         Options must come before dirs",
		"         Dirs are the names/paths of directories to search (cwd by default)",
		"",
		"Options:",
		"        -d        Print loc by directory",
		"             -pd int   Maximum depth of subdirectories to print (default: 1,000)",
		"        -ed str   Directories to exclude (i.e. \"lib,src/utils\")",
		"        -ef str   Files to exclude (i.e. \"index.js,src/main.go\")",
		"        -el str   Languages to exclude (i.e. \"HTML,Plain Text,JSON\")",
		"        -f        Print loc by file",
		"             -mf int   Maximum number of files to print per directory (default: 100,000)",
		"        -id       Include dot directories (excluded by default)",
		"        -il str   Languages to include, all others excluded (i.e. \"Python,JavaScript,C\")",
		"        -ml int   Maximum number of languages to print per directory (default: 1,000)",
		"        -p        Print loc as a percentage of overall total",
		"        -s str    Choose how to sort results [\"loc\", \"size\", \"files\"] (default: \"loc\")",
		"        -sd int   Maximum depth of subdirectories to search (default: 1,000)",
		"",
		"        --help         Print this message and exit",
		"        --license      Print license information and exit",
		"        --version      Print version and exit",
	}

	// Line-by-line text for --license output
	license_message = []string{
		"Source code can be found at github.com/john-bieren/loc",
		"loc is licensed under the MIT license",
		"Copyright (c) 2025 John Bieren",
		"",
		"Languages, extensions, and comment characters sourced from github.com/boyter/scc",
		"scc is licensed under the MIT license",
		"Copyright (c) 2021 Ben Boyter",
		"",
		"",
		"MIT license",
		"",
		"Permission is hereby granted, free of charge, to any person obtaining a copy",
		"of this software and associated documentation files (the \"Software\"), to deal",
		"in the Software without restriction, including without limitation the rights",
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell",
		"copies of the Software, and to permit persons to whom the Software is",
		"furnished to do so, subject to the following conditions:",
		"",
		"The above copyright notice and this permission notice shall be included in all",
		"copies or substantial portions of the Software.",
		"",
		"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR",
		"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,",
		"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE",
		"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER",
		"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,",
		"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE",
		"SOFTWARE.",
	}
)

// Run exit flags, parse list flags, check for invalid inputs
func processFlags() {
	if *version_flag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

	if *license_flag {
		for _, line := range license_message {
			fmt.Println(line)
		}
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

// Custom usage output for --help and relevant error messages
func usage() {
	for _, line := range usage_message {
		fmt.Println(line)
	}
	os.Exit(0)
}
