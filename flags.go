package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

var (
	// printDirFlag is the value of the -d flag.
	printDirFlag = flag.Bool("d", false, "")

	// includeDotDirFlag is the value of the --dot flag.
	includeDotDirFlag = flag.Bool("dot", false, "")

	// excludeDirsFlag is the value of the -ed flag.
	excludeDirsFlag = flag.String("ed", "", "")
	// excludeDirs contains the parsed inputs for the -ed flag.
	excludeDirs []string

	// excludeExtsFlag is the value of the -ee flag.
	excludeExtsFlag = flag.String("ee", "", "")
	// excludeExts contains the parsed inputs for the -ee flag.
	excludeExts []string

	// excludeFilesFlag is the value of the -ef flag.
	excludeFilesFlag = flag.String("ef", "", "")
	// excludeFiles contains the parsed inputs for the -ef flag.
	excludeFiles []string

	// excludeLangsFlag is the value of the -el flag.
	excludeLangsFlag = flag.String("el", "", "")
	// excludeLangs contains the parsed inputs for the -el flag.
	excludeLangs []string

	// printFileFlag is the value of the -f flag.
	printFileFlag = flag.Bool("f", false, "")

	// maxFileReaders is the value of the -fr flag.
	maxFileReaders = flag.Int("fr", runtime.NumCPU(), "")

	// includeDirsFlag is the value of the -id flag.
	includeDirsFlag = flag.String("id", "", "")
	// includeDirs contains the parsed inputs for the -id flag.
	includeDirs []string

	// includeExtsFlag is the value of the -ie flag.
	includeExtsFlag = flag.String("ie", "", "")
	// includeExts contains the parsed inputs for the -ie flag.
	includeExts []string

	// includeFilesFlag is the value of the -if flag.
	includeFilesFlag = flag.String("if", "", "")
	// includeFiles contains the parsed inputs for the -if flag.
	includeFiles []string

	// includeLangsFlag is the value of the -il flag.
	includeLangsFlag = flag.String("il", "", "")
	// includeLangs contains the parsed inputs for the -il flag.
	includeLangs []string

	// maxFilesPrint is the value of the -mf flag.
	maxFilesPrint = flag.Int("mf", 100_000, "")

	// maxTotalsPrint is the value of the -ml flag.
	maxTotalsPrint = flag.Int("ml", 1_000, "")

	// percentagesFlag is the value of the -p flag.
	percentagesFlag = flag.Bool("p", false, "")

	// maxPrintDepth is the value of the -pd flag.
	maxPrintDepth = flag.Int("pd", 1_000, "")

	// suppressWarningsFlag is the value of the -q flag.
	suppressWarningsFlag = flag.Bool("q", false, "")

	// sortColumn is the value of the -s flag.
	sortColumn = flag.String("s", "loc", "")

	// maxSearchDepth is the value of the -sd flag.
	maxSearchDepth = flag.Int("sd", 1_000, "")

	// licenseFlag is the value of the --license flag.
	licenseFlag = flag.Bool("license", false, "")

	// versionFlag is the value of the --version flag.
	versionFlag = flag.Bool("version", false, "")
)

const (
	// licenseMessage is the output of the --license flag.
	licenseMessage = `Source code can be found at github.com/john-bieren/loc
loc is licensed under the MIT license
Copyright (c) 2025 John Bieren

Languages, extensions, and comment characters sourced from github.com/boyter/scc
scc is licensed under the MIT license
Copyright (c) 2021 Ben Boyter

Full text of the MIT license:

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`

	// usageMessage is the output of usage().
	usageMessage = `loc %s
Count lines of code in directories and their subdirectories by language
         Note: multi-line comments and non-comment docstrings are counted as lines of code

Usage: loc [options] [dirs]
         Options must come before dirs
         Option flags cannot be combined (e.g. use -d -f instead of -df)
         Option flags and arguments are case sensitive
         Dirs are the names/paths of directories to search (cwd by default)

Options:
        -d         Print loc by directory
            -pd <int>  Maximum depth of subdirectories to print (default: 1,000)
        --dot      Include dot directories (excluded by default)
        -ed <str>  Directory trees to exclude (name or path, e.g. "node_modules,src/styles")
        -ee <str>  Extensions to exclude (e.g. "json,md,css")
        -ef <str>  Files to exclude (name or path, e.g. "README.md,vendor/htmx.js")
        -el <str>  Languages to exclude (e.g. "HTML,Plain Text,YAML")
        -f         Print loc by file
            -mf <int>  Maximum number of files to print per directory (default: 100,000)
        -fr <int>  Number of file-reading goroutines (default: %d)
        -id <str>  Directory trees to include, excluding others (name or path, e.g. "src,tests/unit")
        -ie <str>  Extensions to include, excluding others (e.g. "go,sh,zig")
        -if <str>  Files to include, excluding others (name or path, e.g. "main.lua,src/index.ts")
        -il <str>  Languages to include, excluding others (e.g. "Python,JavaScript,C++")
        -ml <int>  Maximum number of languages to print per directory (default: 1,000)
        -p         Print loc as a percentage of overall total
        -q         Suppress non-critical error messages
        -s  <str>  How to sort results ["loc", "size", "files"] (default: "loc")
        -sd <int>  Maximum depth of subdirectories to search (default: 1,000)
        --help     Print this message and exit
        --license  Print license information and exit
        --version  Print version and exit`
)

// processFlags runs exit flags, parses string flags, and checks for invalid inputs.
func processFlags() {
	if *versionFlag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

	if *licenseFlag {
		fmt.Println(licenseMessage)
		os.Exit(0)
	}

	if *includeDirsFlag != "" {
		includeDirs = standardizePaths(strings.Split(*includeDirsFlag, ","))
	} else if *excludeDirsFlag != "" {
		excludeDirs = standardizePaths(strings.Split(*excludeDirsFlag, ","))
	}

	if *includeExtsFlag != "" {
		includeExts = strings.Split(*includeExtsFlag, ",")
	} else if *excludeExtsFlag != "" {
		excludeExts = strings.Split(*excludeExtsFlag, ",")
	}

	if *includeFilesFlag != "" {
		includeFiles = standardizePaths(strings.Split(*includeFilesFlag, ","))
	} else if *excludeFilesFlag != "" {
		excludeFiles = standardizePaths(strings.Split(*excludeFilesFlag, ","))
	}

	if *includeLangsFlag != "" {
		includeLangs = strings.Split(*includeLangsFlag, ",")
	} else if *excludeLangsFlag != "" {
		excludeLangs = strings.Split(*excludeLangsFlag, ",")
	}

	if !slices.Contains([]string{"loc", "size", "files"}, *sortColumn) {
		fmt.Printf("-s input \"%s\" is invalid, defaulting to \"loc\"\n", *sortColumn)
		// "loc" is already the default option when sorting results
	}

	if *maxFileReaders < 1 {
		fmt.Printf("-fr input %d is invalid, defaulting to %d\n", *maxFileReaders, runtime.NumCPU())
		*maxFileReaders = runtime.NumCPU()
	}

	if !*printDirFlag {
		*maxPrintDepth = 0
	}
}

/*
standardizePaths formats each item in a slice of path suffixes.
This includes using the proper separators for the user's OS,
handling "." and "..", and ensuring that partial paths have a leading separator
(for example, "-ed lib" would exclude a directory named "my_lib", while changing
"lib" to "/lib" prevents this).
*/
func standardizePaths(input []string) []string {
	var result []string
	for _, path := range input {
		// manually convert leading "." and ".." into their equivalent full paths
		if path[0] == '.' {
			if path[1] == '.' {
				path = strings.Replace(path, "..", parentDir(cwd), 1)
			} else {
				path = strings.Replace(path, ".", cwd, 1)
			}
		}

		// correct separators, handle non-leading "." and "..", and other edge cases
		path = filepath.Clean(path)
		// ensure leading separator
		if !filepath.IsAbs(path) && path[0] != os.PathSeparator {
			path = pathSeparator + path
		}

		result = append(result, path)
	}
	return result
}

// usage is a custom usage output for --help and relevant error messages.
func usage() {
	formattedMessage := fmt.Sprintf(usageMessage, version, *maxFileReaders)
	fmt.Println(formattedMessage)
	os.Exit(0)
}
