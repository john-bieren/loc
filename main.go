// Count the lines of code in a directory
package main

import (
	"flag"
	"fmt"
	"os"
)

// Flag to search and count loc from subdirectories
var include_sub_flag = flag.Bool("a", false, "")

// Flag to print loc by directory
var print_dir_flag = flag.Bool("d", false, "")

// Flag to print loc by file
var print_file_flag = flag.Bool("f", false, "")

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("flags/arguments not properly formatted")
		flag.Usage()
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

	main_dir := newDirectory(dir_path, 0)
	if len(main_dir.children_loc) > 0 {
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
