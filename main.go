//go:generate go run ./generator
//go:generate gofmt -w languages.go

package main

import (
	"flag"
	"fmt"
	"os"
)

// version is the current version of loc.
const version = "v3.1.1-beta"

// cwd is the current working directory.
var cwd string
var totalLoc, totalBytes, totalFiles float64

// main is loc's entry point.
func main() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintln("Error getting cwd:", err))
	}

	// overwrite default usage function to print custom message
	flag.Usage = usage
	flag.Parse()
	processFlags()
	args := flag.Args()

	// dirPaths contains the absolute paths to the directories in the user's arguments.
	var dirPaths []string
	if len(args) == 0 {
		dirPaths = []string{cwd}
	} else {
		// add each directory argument as an absolute path
		for _, path := range args {
			dirPaths = append(dirPaths, toAbsPath(path))
		}
		// make sure each directory is only counted once
		if len(dirPaths) > 1 {
			dirPaths = removeOverlappingDirs(dirPaths)
		}
	}

	// mainDir is the "root" directory from which files and subdirectories are indexed.
	var mainDir *directory
	if len(dirPaths) == 1 {
		mainDir = newDirectory(dirPaths[0], 0)
	} else {
		// increment search depth since this mainDir isn't real but counts as a parent
		*maxSearchDepth++

		// create a fake directory to show totals across multiple directory args
		mainDir = &directory{
			searchSubdirs: true,
			locCounts:     make(map[string]int),
			fileCounts:    make(map[string]int),
			byteCounts:    make(map[string]int),
		}

		for _, path := range dirPaths {
			subdir := newDirectory(path, 1)
			mainDir.subdirectories = append(mainDir.subdirectories, subdir)
		}

		mainDir.countDirLoc()
	}

	if *percentagesFlag {
		totalLoc = float64(sumMapValues(mainDir.locCounts))
		totalBytes = float64(sumMapValues(mainDir.byteCounts))
		totalFiles = float64(sumMapValues(mainDir.fileCounts))
	}

	if len(mainDir.locCounts) > 0 {
		if *printDirFlag {
			mainDir.printTreeLoc()
		} else {
			mainDir.printDirLoc()
			if *printFileFlag {
				mainDir.printFileLoc()
			}
		}
	} else {
		fmt.Println("No code files found")
	}
}
