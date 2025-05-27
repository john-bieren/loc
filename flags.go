package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

var (
	// print_dir_flag is the value of the -d flag.
	print_dir_flag = flag.Bool("d", false, "")

	// include_dot_dir_flag is the value of the --dot flag.
	include_dot_dir_flag = flag.Bool("dot", false, "")

	// exclude_dirs_flag is the value of the -ed flag.
	exclude_dirs_flag = flag.String("ed", "", "")
	// exclude_dirs is a slice of the parsed inputs for the -ed flag.
	exclude_dirs []string

	// exclude_files_flag is the value of the -ef flag.
	exclude_files_flag = flag.String("ef", "", "")
	// exclude_files is a slice of the parsed inputs for the -ef flag.
	exclude_files []string

	// exclude_langs_flag is the value of the -el flag.
	exclude_langs_flag = flag.String("el", "", "")
	// exclude_langs is a slice of the parsed inputs for the -el flag.
	exclude_langs []string

	// print_file_flag is the value of the -f flag.
	print_file_flag = flag.Bool("f", false, "")

	// include_files_flag is the value of the -if flag.
	include_files_flag = flag.String("if", "", "")
	// include_files is a slice of the parsed inputs for the -if flag.
	include_files []string

	// include_langs_flag is the value of the -il flag.
	include_langs_flag = flag.String("il", "", "")
	// include_langs is a slice of the parsed inputs for the -il flag.
	include_langs []string

	// max_print_files is the value of the -mf flag.
	max_print_files = flag.Int("mf", 100_000, "")

	// max_print_totals is  the value of the -ml flag.
	max_print_totals = flag.Int("ml", 1_000, "")

	// percentages_flag is the value of the -p flag.
	percentages_flag = flag.Bool("p", false, "")

	// max_print_depth is the value of the -pd flag.
	max_print_depth = flag.Int("pd", 1_000, "")

	// sort_column is the value of the -s flag.
	sort_column = flag.String("s", "loc", "")

	// max_search_depth is the value of the -sd flag.
	max_search_depth = flag.Int("sd", 1_000, "")

	// suppress_warnings is the value of the -q flag.
	suppress_warnings = flag.Bool("q", false, "")

	// license_flag is the value of the --license flag.
	license_flag = flag.Bool("license", false, "")

	// version_flag is the value of the --version flag.
	version_flag = flag.Bool("version", false, "")

	// usage_message is the line-by-line text for usage function output.
	usage_message = []string{
		fmt.Sprintf("loc %s", version),
		"Count lines of code in directories and their subdirectories by language",
		"         Note: multi-line comments and non-comment docstrings are counted as lines of code",
		"",
		"Usage: loc [options] [dirs]",
		"         Options must come before dirs",
		"         Option flags cannot be combined (e.g. use -d -f instead of -df)",
		"         Dirs are the names/paths of directories to search (cwd by default)",
		"",
		"Options:",
		"        -d        Print loc by directory",
		"            -pd int   Maximum depth of subdirectories to print (default: 1,000)",
		"        --dot     Include dot directories (excluded by default)",
		"        -ed str   Directories to exclude (name or path, e.g. \"lib,src/utils\")",
		"        -ef str   Files to exclude (name or path, e.g. \"index.js,src/main.go\")",
		"        -el str   Languages to exclude (e.g. \"HTML,Plain Text,JSON\")",
		"        -f        Print loc by file",
		"            -mf int   Maximum number of files to print per directory (default: 100,000)",
		"        -if str   Files to include (name or path, e.g. \"main.py,src/main.c\")",
		"        -il str   Languages to include, all others excluded (e.g. \"Python,JavaScript,C\")",
		"        -ml int   Maximum number of languages to print per directory (default: 1,000)",
		"        -p        Print loc as a percentage of overall total",
		"        -s  str   Choose how to sort results [\"loc\", \"size\", \"files\"] (default: \"loc\")",
		"        -sd int   Maximum depth of subdirectories to search (default: 1,000)",
		"        -q        Suppress non-critical error messages",
		"        --help    Print this message and exit",
		"        --license   Print license information and exit",
		"        --version   Print version and exit",
	}

	// license_message is the line-by-line text for --license output.
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

// processFlags runs exit flags, parses string flags, and checks for invalid inputs.
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
		exclude_dirs = standardizeSeparators(strings.Split(*exclude_dirs_flag, ","))
	}
	if *exclude_files_flag != "" {
		exclude_files = standardizeSeparators(strings.Split(*exclude_files_flag, ","))
	}
	if *exclude_langs_flag != "" {
		exclude_langs = strings.Split(*exclude_langs_flag, ",")
	}
	if *include_files_flag != "" {
		include_files = standardizeSeparators(strings.Split(*include_files_flag, ","))
	}
	if *include_langs_flag != "" {
		include_langs = strings.Split(*include_langs_flag, ",")
	}

	if !slices.Contains([]string{"loc", "size", "files"}, *sort_column) {
		fmt.Printf("-s input \"%s\" is invalid, defaulting to \"loc\"\n", *sort_column)
	}
}

/*
standardizeSeparators corrects path separators in a slice of paths.
This includes using the proper separators for the user's OS, and ensuring that
there is a leading separator, as these paths will match to entries which
contain them as a suffix. For example, "-ed lib" would exclude a directory
named "my_lib", while changing "lib" to "/lib" prevents this.
*/
func standardizeSeparators(input []string) []string {
	var result []string
	windows := string(os.PathSeparator) == "\\"
	for _, path := range input {
		if windows {
			path = strings.ReplaceAll(path, "/", "\\")
			path = strings.Trim(path, "\\")
			path = fmt.Sprintf("\\%s", path)
			result = append(result, path)
		} else {
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.Trim(path, "/")
			path = fmt.Sprintf("/%s", path)
			result = append(result, path)
		}
	}
	return result
}

// usage is a custom usage output for --help and relevant error messages.
func usage() {
	for _, line := range usage_message {
		fmt.Println(line)
	}
	os.Exit(0)
}
