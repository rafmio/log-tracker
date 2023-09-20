package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	// t := time.Now()
	logstring := "Aug 20 00:01:45 localhost kernel: [3217949.737464] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=109.205.213.90 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=240 ID=57645 PROTO=TCP SPT=51638 DPT=2106 WINDOW=1024 RES=0x00 SYN URGP=0"
	logtime := logstring[0:12]
	year := strconv.Itoa(time.Now().Year())

	fmt.Println("Year:", year)
	logtime = logtime + " " + year
	fmt.Println(logtime)

	tparsed, err := time.Parse("Jan 2 15:04 2006", logtime)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(tparsed)
	fmt.Println(logtime)
}
