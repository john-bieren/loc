package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	fullPath string
	name     string
	fileType string
	bytes    int
	loc      int
}

// countFileLoc counts the lines of code in f.
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
func newFile(path, lang string, size int64) *file {
	self := &file{
		fullPath: path,
		name:     filepath.Base(path),
		fileType: lang,
		bytes:    int(size),
	}
	self.countFileLoc()
	return self
}
