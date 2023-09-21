package main

import (
	"bufio"
	"os"
)

func ParseSrcData(srcData *os.File, logs *[]Log) error {
	// create a scanner for text processing
	scanner := bufio.NewScanner(srcData)
	for scanner.Scan() {
		// parse line func
	}
}
