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

// countFileLoc counts the lines of code in a file.
func (f *file) countFileLoc() {
	file, err := os.Open(f.full_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	com_chars, has_comments := single_line_comment_chars[f.file_type]
	reader := bufio.NewReader(file)
	var end_of_file, skip_line bool
	for !end_of_file {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				end_of_file = true
			} else {
				fmt.Println("Error reading line:", err)
				continue
			}
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if has_comments {
			for _, char := range com_chars {
				if strings.HasPrefix(line, char) {
					skip_line = true
					break
				}
			}
			if skip_line {
				skip_line = false
				continue
			}
		}

		f.loc++
	}
}

// newFile is the constructor for instances of the file struct.
func newFile(path string, dir_parents int, size int64) *file {
	self := &file{
		full_path: path,
		rel_path:  relPath(path, dir_parents),
		name:      filepath.Base(path),
		bytes:     int(size),
	}
	self.file_type, self.is_code = filenames[self.name]
	if !self.is_code {
		self.file_type, self.is_code = extensions[filepath.Ext(path)]
	}

	if len(include_langs) > 0 {
		self.is_code = false
		for _, incl := range include_langs {
			if self.file_type == incl {
				self.is_code = true
			}
		}
	} else {
		for _, excl := range exclude_langs {
			if self.file_type == excl {
				self.is_code = false
			}
		}
	}

	if self.is_code {
		self.countFileLoc()
	}
	return self
}
