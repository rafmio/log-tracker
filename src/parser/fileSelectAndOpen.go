package parser

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// var mapFiles = make(map[string]time.Time)

var ErrGetStatInfo = errors.New("can't get file info via Stat()")

// Searches for files with a name corresponding to the pattern, selects
// the non-zero and most recent one, and then opens it and returns a pointer
// to the open file.
// The target folder is passed to the function from the environment variable
func SelectAndOpen(directory string) (*os.File, error) {
	// create variable for storing filenames and time
	var mapFiles = make(map[string]time.Time)

	// looking at the filenames in the entire directory:
	files, err := filepath.Glob(filepath.Join(directory, "ufw.log*"))
	if err != nil {
		return nil, err
	}

	// fills the map with the names of files corresponding to the pattern
	for _, filename := range files {
		fi, err := os.Stat(filename) // getting information about the file
		if err != nil {
			return nil, err
		}

		// exclude archives and empty files
		if fi.Size() == 0 || strings.Contains(filename, ".gz") {
			continue
		}
		mapFiles[filename] = fi.ModTime()
	}

	// find the latest file
	latestTime := mapFiles[files[0]]
	var latestFile string
	for fileName, tm := range mapFiles {
		if latestTime.Before(tm) {
			latestTime = tm
			latestFile = fileName
		}
	}

	file, err := os.Open(latestFile)
	if err != nil {
		return nil, err
	}

	return file, nil
}
