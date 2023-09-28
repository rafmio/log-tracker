package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// var SrcFileName string = "ufw-1-line.log"
var SrcFileName string = "ufw.log"

// set log-file path and set golbal log-file var for the application
var logFilePath string = "logfile.log"
var LogFile *os.File

// set log-file settings
func init() {
	LogFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("ERROR: opening log-file")
	} else {
		log.SetOutput(LogFile)
		log.SetFlags(log.LstdFlags)
		log.Println("INFO: log file has been opened")
	}
}

type Log struct {
	Date      time.Time
	SrcIP     net.IP
	PacketLen int
	Ttl       int
	PacketId  int
	SrcPort   int
	DptPort   int
	Window    int
}

func main() {
	defer LogFile.Close()

	// create slice for store parsed log structs
	logs := make([]Log, 0)

	// open source file
	srcData, err := os.OpenFile(SrcFileName, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal("ERROR: opening src file")
	} else {
		log.Println("INFO: the src file has been read")
	}

	// pass the read data and tmp storage for parsing
	err = ParseSrcData(srcData, &logs)
	if err != nil {
		fmt.Println("ERROR: main()", err.Error())
	}

	fmt.Println(
		logs[32].Date,
		logs[32].SrcIP.String(),
		logs[32].PacketLen,
		logs[32].Ttl,
		logs[32].PacketId,
		logs[32].SrcPort,
		logs[32].DptPort,
		logs[32].Window,
	)
	fmt.Println("The length of the logs is:", len(logs))

	srcData.Close()
}
