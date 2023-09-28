package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseSrcData(srcData *os.File, logs *[]Log) error {
	var err error = nil

	// create a scanner for text processing
	scanner := bufio.NewScanner(srcData)
	for scanner.Scan() {
		// declaring the log structure to fill in the fields
		var logLine Log

		// saving the log-line to a variable
		logLineStr := scanner.Text()

		if err := scanner.Err(); err != nil {
			log.Println("ERROR: scanner", err.Error())
			return err
		}

		// splitting the string into tokens
		tokens := strings.Fields(logLineStr)

		// tokens[0] is a month and tokens[1] is a day. Add year
		dateStr := tokens[1] + " " + tokens[0] + " " + strconv.Itoa(time.Now().Year())
		// tokens[2] is a time, add it
		dateStr = dateStr + " " + tokens[2]

		// initialize the Date field of the log structure by parsed date
		date, err := time.Parse("02 Jan 2006 15:04:05", dateStr)
		if err != nil {
			log.Println("ERROR: time.Parse(dateStr)", err.Error())
			return err
		} else {
			logLine.Date = date
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

		// filling in the structure fields
		logLine.SrcIP = net.ParseIP(prefixesMap["SRC"])

		if prefixesMap["LEN"] == "" {
			logLine.PacketLen = 0
			log.Println("WARNING: the LEN is empty")
		} else {
			logLine.PacketLen, err = strconv.Atoi(prefixesMap["LEN"])
			if err != nil {
				log.Println("ERROR: converting string LEN to int")
				return err
			}
		}

		if prefixesMap["TTL"] == "" {
			logLine.Ttl = 0
			log.Println("WARNING: the TTL is empty")
		} else {
			logLine.Ttl, err = strconv.Atoi(prefixesMap["TTL"])
			if err != nil {
				log.Println("ERROR: converting string TTL to int")
				return err
			}
		}

		if prefixesMap["ID"] == "" {
			logLine.PacketId = 0
			log.Println("WARINING: the ID is empty")
		} else {
			logLine.PacketId, err = strconv.Atoi(prefixesMap["ID"])
			if err != nil {
				fmt.Println("ERROR: converting string ID to int")
				return err
			}
		}

		if prefixesMap["SPT"] == "" {
			logLine.SrcPort = 0
			log.Println("WARNING: the SPT is empty")
		} else {
			logLine.SrcPort, err = strconv.Atoi(prefixesMap["SPT"])
			if err != nil {
				fmt.Println("ERROR: converting string SPT to int")
				return err
			}
		}

		if prefixesMap["DPT"] == "" {
			logLine.DptPort = 0
			log.Println("WARNING: the DPT is empty")
		} else {
			logLine.DptPort, err = strconv.Atoi(prefixesMap["DPT"])
			if err != nil {
				fmt.Println("ERROR: converting string DPT to int")
				return err
			}
		}

		if prefixesMap["WINDOW"] == "" {
			logLine.Window = 0
			log.Println("WARNING: the WINDOW is empty")
		} else {
			logLine.Window, err = strconv.Atoi(prefixesMap["WINDOW"])
			if err != nil {
				fmt.Println("ERROR: converting string WINDOW to int")
				return err
			}
		}

		// adding the structure to the slice []Log
		*logs = append(*logs, logLine)
	}

	log.Println("INFO: parsing is done")
	return err
}
