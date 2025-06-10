package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// treeFileHeadersPrinted tracks whether file column headers have been printed by printTreeLoc.
// These must be printed once but will not print on the first directory if it has no files.
var treeFileHeadersPrinted bool

type directory struct {
	fullPath       string
	name           string
	parents        int
	searchSubdirs  bool
	subdirectories []*directory
	files          []*file
	locCounts      map[string]int
	fileCounts     map[string]int
	byteCounts     map[string]int
}

// searchDir indexes the directory's files and subdirectories.
func (d *directory) searchDir() {
	entries, err := os.ReadDir(d.fullPath)
	if err != nil {
		warn("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		entryName := entry.Name()
		fullPath := filepath.Join(d.fullPath, entryName)
		info, err := os.Stat(fullPath)
		if err != nil {
			// specify errors from inaccessible entries, a common case
			if os.IsNotExist(err) {
				warn("Cannot access directory entry:", err)
			} else {
				warn("Error checking directory entry:", err)
			}
			continue
		}

		if info.IsDir() {
			if d.searchSubdirs {
				if !*includeDotDirFlag && strings.HasPrefix(entryName, ".") {
					continue
				}

				var skipDir bool
				for _, excl := range excludeDirs {
					if strings.HasSuffix(fullPath, excl) {
						skipDir = true
						break
					}
				}
				if skipDir {
					continue
				}

				subdir := newDirectory(fullPath, d.parents+1)
				d.subdirectories = append(d.subdirectories, subdir)
			}
		} else {
			var skipFile bool
			// check for matches with included/excluded files
			if len(includeFiles) > 0 {
				skipFile = true
				for _, incl := range includeFiles {
					if strings.HasSuffix(fullPath, incl) {
						skipFile = false
						break
					}
				}
			} else {
				for _, excl := range excludeFiles {
					if strings.HasSuffix(fullPath, excl) {
						skipFile = true
						break
					}
				}
			}
			if skipFile {
				continue
			}

			// determine file's language by file name
			fileType, isCode := filenames[entryName]
			if !isCode {
				// determine file's language by file extension
				fileType, isCode = extensions[filepath.Ext(entryName)]
			}
			if !isCode {
				continue
			}

			// check for matches with included/excluded languages
			if len(includeLangs) > 0 {
				for _, incl := range includeLangs {
					if fileType == incl {
						skipFile = true
					}
				}
			} else {
				for _, excl := range excludeLangs {
					if fileType == excl {
						skipFile = false
					}
				}
			}
			if skipFile {
				continue
			}

			size := info.Size()
			file := newFile(fullPath, fileType, d.parents, size)
			d.files = append(d.files, file)
		}
	}
}

// countDirLoc counts the lines of code for each language in all indexed files.
func (d *directory) countDirLoc() {
	for _, file := range d.files {
		d.locCounts[file.fileType] += file.loc
		d.fileCounts[file.fileType]++
		d.byteCounts[file.fileType] += file.bytes
	}

	if d.searchSubdirs {
		for _, subdir := range d.subdirectories {
			for fileType, loc := range subdir.locCounts {
				d.locCounts[fileType] += loc
			}
			for fileType, n := range subdir.fileCounts {
				d.fileCounts[fileType] += n
			}
			for fileType, b := range subdir.byteCounts {
				d.byteCounts[fileType] += b
			}
		}
	}
}

// printDirLoc prints loc by file type for the directory.
func (d *directory) printDirLoc() {
	if len(d.locCounts) == 0 {
		return
	}
	indent := strings.Repeat("    ", d.parents)

	// print directory name, if applicable
	if *printDirFlag && d.parents > 0 {
		fmt.Printf("%s%s/\n", indent, d.name)
		indent += " " // loc totals should have an extra space if directory names are printed
	}

	// print column labels on first directory
	if d.parents == 0 {
		fmt.Printf("\033[1m%sLanguage: loc | size | files\033[0m\n", indent)
	}

	// print loc total if multiple languages are present
	if len(d.locCounts) > 1 {
		if *percentagesFlag && d.parents > 0 {
			fmt.Printf(
				"%s%d langs: %.1f%% | %.1f%% | %.1f%%\n",
				indent, len(d.locCounts),
				float64(sumMapValues(d.locCounts))/totalLoc*100,
				float64(sumMapValues(d.byteCounts))/totalBytes*100,
				float64(sumMapValues(d.fileCounts))/totalFiles*100,
			)
		} else {
			fmt.Printf(
				"%s%d langs: %s | %s | %s\n",
				indent, len(d.locCounts),
				addCommas(sumMapValues(d.locCounts)),
				formatByteCount(sumMapValues(d.byteCounts)),
				addCommas(sumMapValues(d.fileCounts)),
			)
		}
	}

	// keys contains the file type keys sorted by their sortColumn values.
	var keys []string
	switch *sortColumn {
	case "size":
		keys = sortKeys(d.byteCounts)
	case "files":
		keys = sortKeys(d.fileCounts)
	default:
		keys = sortKeys(d.locCounts)
	}
	// print loc totals by file type
	for i, fileType := range keys {
		// print language total even if -ml=0 if there's only one language
		if i+1 > *maxTotalsPrint && len(d.locCounts) > 1 {
			break
		}
		if *percentagesFlag && !(len(d.locCounts) == 1 && d.parents == 0) {
			fmt.Printf(
				"%s%s: %.1f%% | %.1f%% | %.1f%%\n",
				indent, fileType,
				float64(d.locCounts[fileType])/totalLoc*100,
				float64(d.byteCounts[fileType])/totalBytes*100,
				float64(d.fileCounts[fileType])/totalFiles*100,
			)
		} else {
			fmt.Printf(
				"%s%s: %s | %s | %s\n",
				indent, fileType,
				addCommas(d.locCounts[fileType]),
				formatByteCount(d.byteCounts[fileType]),
				addCommas(d.fileCounts[fileType]),
			)
		}
	}
}

// printTreeLoc prints loc by file type for the directory and its subdirectories, includes files if -f used.
func (d *directory) printTreeLoc() {
	d.printDirLoc()

	if *printFileFlag {
		indent := strings.Repeat("    ", d.parents+1)
		if !treeFileHeadersPrinted && len(d.files) > 0 {
			fmt.Printf("\033[1m%sloc | size - file\033[0m\n", indent)
			treeFileHeadersPrinted = true
		}

		for i, file := range sortFiles(d.files, *sortColumn) {
			if i+1 > *maxFilesPrint {
				break
			}
			if *percentagesFlag {
				fmt.Printf(
					"%s%.1f%% | %.1f%% - %s\n",
					indent,
					float64(file.loc)/totalLoc*100,
					float64(file.bytes)/totalBytes*100,
					file.name,
				)
			} else {
				fmt.Printf(
					"%s%s | %s - %s\n",
					indent,
					addCommas(file.loc),
					formatByteCount(file.bytes),
					file.name,
				)
			}
		}
	}

	if d.searchSubdirs && *maxPrintDepth >= d.parents+1 {
		// sort the subdirectories by the selected sort column
		sort.Slice(d.subdirectories, func(i, j int) bool {
			switch *sortColumn {
			case "size":
				return sumMapValues(d.subdirectories[i].byteCounts) > sumMapValues(d.subdirectories[j].byteCounts)
			case "files":
				return sumMapValues(d.subdirectories[i].fileCounts) > sumMapValues(d.subdirectories[j].fileCounts)
			default:
				return sumMapValues(d.subdirectories[i].locCounts) > sumMapValues(d.subdirectories[j].locCounts)
			}
		})

		for _, subdir := range d.subdirectories {
			subdir.printTreeLoc()
		}
	}
}

// printFileLoc prints loc by file for all files counted.
func (d *directory) printFileLoc() {
	// files contains the directory's files and is used to sort them.
	var files []*file
	files = d.appendFiles(files)

	fmt.Println("\033[1m loc | size - file\033[0m")
	for i, file := range sortFiles(files, *sortColumn) {
		if i+1 > *maxFilesPrint {
			break
		}
		if *percentagesFlag {
			fmt.Printf(
				" %.1f%% | %.1f%% - %s\n",
				float64(file.loc)/totalLoc*100,
				float64(file.bytes)/totalBytes*100,
				file.relPath,
			)
		} else {
			fmt.Printf(
				" %s | %s - %s\n",
				addCommas(file.loc),
				formatByteCount(file.bytes),
				file.relPath,
			)
		}
	}
}

// appendFiles appends the contents of d.files to the input slice.
func (d *directory) appendFiles(input []*file) []*file {
	input = append(input, d.files...)
	if d.searchSubdirs {
		for _, subdir := range d.subdirectories {
			input = subdir.appendFiles(input)
		}
	}
	return input
}

// newDirectory is the constructor for instances of the directory struct.
func newDirectory(path string, parents int) *directory {
	self := &directory{
		fullPath:      path,
		name:          filepath.Base(path),
		parents:       parents,
		searchSubdirs: parents+1 <= *maxSearchDepth,
		locCounts:     make(map[string]int),
		fileCounts:    make(map[string]int),
		byteCounts:    make(map[string]int),
	}
	self.searchDir()
	self.countDirLoc()
	return self
}
