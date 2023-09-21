package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseSrcData(srcData *os.File, logs *[]Log) {

	// create a scanner for text processing
	scanner := bufio.NewScanner(srcData)
	for scanner.Scan() {
		// declaring variable of type Log struct for appending to []Log
		// var log Log

		// saving the log-line to a variable
		logLine := scanner.Text()

		// splitting the string into tokens
		tokens := strings.Fields(logLine)

		// tokens[0] is a month and tokens[1] is a day. Add year
		dateStr := tokens[1] + " " + tokens[0] + " " + strconv.Itoa(time.Now().Year())

		// initialize the Date field of the log structure by parsed date
		date, err := time.Parse("02 Jan 2006", dateStr)
		if err != nil {
			fmt.Println("time.Parse(dateStr):", err.Error())
		}

		// tokens[2] is time. Save it to log.HMtime
		hmTime, err := time.Parse("15:04", tokens[2])
		if err != nil {
			fmt.Println("time.Parse(tokens[2])", err.Error)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err.Error())
		}

		logs = append(logs, Log{Date: date, HMtime: hmTime})

	}

}
