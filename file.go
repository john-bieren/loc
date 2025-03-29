package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	full_path string
	rel_path  string
	name      string
	file_type string
	is_code   bool
	loc       int
	bytes     int
}

// Count lines of code in a file
func (f *file) countFileLoc() {
	file, err := os.Open(f.full_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	end_of_file := false

	for !end_of_file {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				end_of_file = true
			} else {
				fmt.Println("Error reading line:", err)
				break
			}
		}
		if strings.TrimSpace(line) != "" {
			f.loc++
		}
	}
}

// Constructor for instances of file struct
func newFile(path string, dir_parents int, size int64) *file {
	self := &file{
		full_path: path,
		rel_path:  relPath(path, dir_parents),
		name:      filepath.Base(path),
		bytes:     int(size),
	}
	self.file_type, self.is_code = file_languages[self.name]
	if !self.is_code {
		self.file_type, self.is_code = langauges[filepath.Ext(path)]
	}
	if self.is_code {
		self.countFileLoc()
	}
	return self
}
