package parser

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrGetStatInfo = errors.New("can't get file info via Stat()")

// Searches for files with a name corresponding to the pattern, selects
// the non-zero and most recent one, and then opens it and returns a pointer
// to the open file.
// The target folder is passed to the function from the environment variable
func SelectAndOpen(fileConfig FileConfig) (*os.File, error) {
	// check if FileConfig's filelds are empty
	if fileConfig.Pattern == "" || fileConfig.Directory == "" || fileConfig.ExcludePattern == "" {
		return nil, errors.New("fileConfig's fields are empty")
	}

	// check if the directory exists
	if _, err := os.Stat(fileConfig.Directory); os.IsNotExist(err) {
		return nil, errors.New("directory doesn't exist")
	}

	// create variable for storing filenames and time
	var mapFiles = make(map[string]time.Time)

	// looking at the filenames in the entire directory:
	files, err := filepath.Glob(filepath.Join(fileConfig.Directory, fileConfig.Pattern))
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
		if fi.Size() == 0 || strings.Contains(filename, fileConfig.ExcludePattern) {
			continue
		}
		mapFiles[filename] = fi.ModTime()
	}

	// var latestFile string // variable for storing latest file name
	latestFile := files[0]
	latestTime := mapFiles[files[0]]

	// find the latest file
	if len(mapFiles) > 1 {
		for fileName, tm := range mapFiles {
			if latestTime.Before(tm) {
				latestTime = tm
				latestFile = fileName
			}
		}
	}

	file, err := os.Open(latestFile)
	if err != nil {
		log.Println("SelectAndOpen() - opening log-file:", err)
		return nil, err
	}

	return file, nil
}
