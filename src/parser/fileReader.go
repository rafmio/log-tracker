package parser

import (
	"bufio"
	"fmt"
	"os"
)

func FileReader(file *os.File, filePosition int) []string {

	slsStr := make([]string, 0)

	_, err := file.Seek(int64(filePosition), os.SEEK_SET)
	if err != nil {
		fmt.Println("setting file position:", err.Error())
		// TODO: log error
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slsStr = append(slsStr, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanning lines:", err.Error())
		// TODO: log error
	}

	return slsStr
}
