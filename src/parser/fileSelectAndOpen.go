package parser

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var mapFiles = make(map[string]time.Time)

var ErrGetStatInfo = errors.New("can't get file info via Stat()")

// pass target directory via env
func SelectAndOpen(directory string) (*os.File, error) {
	// looking at the filenames in the entire directory:
	files, err := filepath.Glob(filepath.Join(directory, "ufw.log*"))
	if err != nil {
		return nil, err
	}

	// fill mapFiles, check if file is empty
	for _, filename := range files {

		fi, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

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
