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

func ParseSrcData(srcData *os.File, logs *[]Log) {

	// create a scanner for text processing
	scanner := bufio.NewScanner(srcData)
	for scanner.Scan() {
		// declaring the log structure to fill in the fields
		var log Log

		// saving the log-line to a variable
		logLine := scanner.Text()

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
		} else {
			log.Date = date
		}

		// tokens[11] is the IP-address
		if strings.Index(tokens[11], "SRC=") != -1 { // check if tokens[10] contains "SRC="
			ip := strings.TrimPrefix(tokens[11], "SRC=") // if true - trim it
			log.SrcIP = net.ParseIP(ip)                  // initialize SrcIP field with ip converted to net.IP
		} else {
			fmt.Println("Can't parse IP")
		}

		// tokens[13] is the length of the packet
		if strings.Index(tokens[13], "LEN=") != -1 { // check if tokens[13] contiains "LEN="
			len := strings.TrimPrefix(tokens[13], "LEN=") // if true - trim it
			log.PacketLen, err = strconv.Atoi(len)        // initialize PacketLen with it converted to int
			if err != nil {
				fmt.Println("Can't convert string to int", err.Error())
			}
		} else {
			fmt.Println("Can't parse LEN=")
		}

		// tokens[16] is the packet's TTL (Ttl)
		if strings.Index(tokens[16], "TTL=") != -1 { // check if tokens[16] contains "TTL="
			ttl := strings.TrimPrefix(tokens[16], "TTL=") // if true - trim it
			log.Ttl, err = strconv.Atoi(ttl)              // initialize Ttl with converted to int
			if err != nil {
				fmt.Println("Can't convert string go int", err.Error())
			}
		} else {
			fmt.Println("Can't parse TTL=")
		}

		// tokens[17] is the packet's ID (PacketID)
		if strings.Index(tokens[17], "ID=") != -1 {
			id := strings.TrimPrefix(tokens[17], "ID=")
			log.PacketId, err = strconv.Atoi(id)
			if err != nil {
				fmt.Println("Can't convert string to int")
			}
		} else {
			fmt.Println("Can't parse ID=")
		}

		// tokens[19] is the source port (SrcPort)
		if strings.Index(tokens[19], "SPT=") != -1 {
			spt := strings.TrimPrefix(tokens[19], "SPT=")
			log.SrcPort, err = strconv.Atoi(spt)
			if err != nil {
				fmt.Println("Can't convert string to int")
			}
		} else {
			fmt.Println("Can't parse SPT=")
		}

		if strings.Index(tokens[20], "DPT=") != -1 {
			dst := strings.TrimPrefix(tokens[20], "DPT=")
			log.DstPort, err = strconv.Atoi(dst)
			if err != nil {
				fmt.Println("Can't convert string to int")
			}
		} else {
			fmt.Println("Can't parse DPT=")
		}

		if strings.Index(tokens[21], "WINDOW=") != -1 {
			window := strings.TrimPrefix(tokens[21], "WINDOW=")
			log.Window, err = strconv.Atoi(window)
			if err != nil {
				fmt.Println("Can't convert string to int")
			}
		} else {
			fmt.Println("Can't parse WINDOW=")
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err.Error())
		}

		*logs = append(*logs, log)
	}
}

/*
TODO: Combine the monotonous parsing operations of each token
into a separate function for better readability and convenient maintenance
*/
