package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	fullPath string
	relPath  string
	name     string
	fileType string
	isCode   bool
	loc      int
	bytes    int
}

// countFileLoc counts the lines of code in a file.
func (f *file) countFileLoc() {
	file, err := os.Open(f.fullPath)
	if err != nil {
		warn("Error opening file:", err)
		return
	}
	defer file.Close()

	comChars, hasComments := singleLineCommentChars[f.fileType]
	reader := bufio.NewReader(file)
	var endOfFile, skipLine bool
	for !endOfFile {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				endOfFile = true
			} else {
				warn("Error reading line:", err)
				continue
			}
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if hasComments {
			for _, char := range comChars {
				if strings.HasPrefix(line, char) {
					skipLine = true
					break
				}
			}
			if skipLine {
				skipLine = false
				continue
			}
		}

		f.loc++
	}
}

// newFile is the constructor for instances of the file struct.
func newFile(path string, dirParents int, size int64) *file {
	self := &file{
		fullPath: path,
		relPath:  relPath(path, dirParents),
		name:     filepath.Base(path),
		bytes:    int(size),
	}
	self.fileType, self.isCode = filenames[self.name]
	if !self.isCode {
		self.fileType, self.isCode = extensions[filepath.Ext(path)]
	}

	if len(includeLangs) > 0 {
		self.isCode = false
		for _, incl := range includeLangs {
			if self.fileType == incl {
				self.isCode = true
			}
		}
	} else {
		for _, excl := range excludeLangs {
			if self.fileType == excl {
				self.isCode = false
			}
		}
	}

	if self.isCode {
		self.countFileLoc()
	}
	return self
}
