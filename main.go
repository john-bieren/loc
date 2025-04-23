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
	if len(args) > 1 {
		fmt.Println("flags/arguments not properly formatted")
		flag.Usage()
	}
	if *version_flag {
		fmt.Println("loc", version)
		os.Exit(0)
	}

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
