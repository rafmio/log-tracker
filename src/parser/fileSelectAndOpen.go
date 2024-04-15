package parser

import (
	"os"
	"path/filepath"
	"time"
)

// pass target directory via env
func SelectAndOpen(directory string) (*os.File, error) {
	// looking at the filenames in the entire directory:
	files, err := filepath.Glob(filepath.Join(directory, "ufw.log*"))
	if err != nil {
		return nil, err
	}

	// slice for storing creating times
	creationTimes := make([]time.Time, 0, len(files))

	// fill the slice with creation time values
	for _, filename := range files {
		info, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}
		creationTimes = append(creationTimes, info.ModTime())
	}

	latestTime := creationTimes[0] // start position for range
	latestFile := files[0]

	for i, tm := range creationTimes {
		if latestTime.Before(tm) {
			latestTime = tm
			latestFile = files[i]
		}
	}

	file, err := os.Open(latestFile)
	if err != nil {
		return nil, err
	}

	return file, nil
}
