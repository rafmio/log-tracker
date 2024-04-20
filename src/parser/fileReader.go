package parser

import (
	"bufio"
	"io"
	"log"
	"os"
)

func FileReader(file *os.File, filePosition int64) ([]string, error) {

	slsStr := make([]string, 0)

	_, err := file.Seek(filePosition, io.SeekStart)
	if err != nil {
		log.Println("os.Seek():", err)
		return slsStr, err
	}

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
