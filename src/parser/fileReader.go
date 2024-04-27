package parser

import (
	"bufio"
	"io"
	"log"
	"os"
)

// The function takes two arguments: a pointer to an open file and
// the value of the file position from which to start reading the file.
func FileReader(file *os.File, filePosition int64) ([]string, error) {

	slsStr := make([]string, 0)

	// setting the file position
	_, err := file.Seek(filePosition, io.SeekStart)
	if err != nil {
		log.Println("os.Seek():", err)
		return slsStr, err
	}

	// read the file line by line and write each line separately into a slice
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slsStr = append(slsStr, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("scanning lines:", err.Error())
	}

	return slsStr, nil
}
