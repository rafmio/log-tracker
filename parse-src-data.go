package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseSrcData(srcData *os.File, logs *[]Log) error {
	var counter int = 0
	var err error = nil
	// create a scanner for text processing
	scanner := bufio.NewScanner(srcData)
	for scanner.Scan() {
		// declaring the log structure to fill in the fields
		var log Log

		// saving the log-line to a variable
		logLine := scanner.Text()

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err.Error())
			return err
		}

		// splitting the string into tokens
		tokens := strings.Fields(logLine)

		// tokens[0] is a month and tokens[1] is a day. Add year
		dateStr := tokens[1] + " " + tokens[0] + " " + strconv.Itoa(time.Now().Year())
		// tokens[2] is a time, add it
		dateStr = dateStr + " " + tokens[2]

		// initialize the Date field of the log structure by parsed date
		date, err := time.Parse("02 Jan 2006 15:04:05", dateStr)
		if err != nil {
			fmt.Println("time.Parse(dateStr):", err.Error())
			return err
		} else {
			log.Date = date
		}

		// declare a map of prefixes
		prefixesMap := map[string]string{
			"SRC":    "",
			"LEN":    "",
			"TTL":    "",
			"ID":     "",
			"SPT":    "",
			"DPT":    "",
			"WINDOW": "",
		}

		// iterate over the slice starting from the 12th element
		for i := 10; i < len(tokens); i++ {
			// divide an element of type "LEN=40"
			// into two elements with a separator "="
			prefixValue := strings.Split(tokens[i], "=")
			_, ok := prefixesMap[prefixValue[0]] // check if the key exists
			if ok {
				// if key exists - fill appropriate value
				prefixesMap[prefixValue[0]] = prefixValue[1]
			} else {
				continue
			}
		}

		log.SrcIP = net.ParseIP(prefixesMap["SRC"])

		log.PacketLen, err = strconv.Atoi(prefixesMap["LEN"])
		if err != nil {
			fmt.Println("converting string LEN to int")
			fmt.Println("counter: ", counter)
			return err
		}

		log.Ttl, err = strconv.Atoi(prefixesMap["TTL"])
		if err != nil {
			fmt.Println("converting string TTL to int")
			return err
		}

		log.PacketId, err = strconv.Atoi(prefixesMap["ID"])
		if err != nil {
			fmt.Println("converting string ID to int")
			return err
		}

		log.SrcPort, err = strconv.Atoi(prefixesMap["SPT"])
		if err != nil {
			fmt.Println("converting string SPT to int")
			return err
		}

		log.DptPort, err = strconv.Atoi(prefixesMap["DPT"])
		if err != nil {
			fmt.Println("converting string DPT to int")
			return err
		}

		if prefixesMap["WINDOW"] == "" {
			log.Window = 0
		} else {
			log.Window, err = strconv.Atoi(prefixesMap["WINDOW"])
			if err != nil {
				fmt.Println("converting string WINDOW to int")
				return err
			}
		}

		*logs = append(*logs, log)
	}

	counter++
	return err
}
