package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// var SrcFileName string = "ufw-1-line.log"
var SrcFileName string = "ufw.log"

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
	logs := make([]Log, 0) // create slice for store parsed log structs

	srcData, err := os.OpenFile(SrcFileName, os.O_RDONLY, 0) // open source file
	if err != nil {
		fmt.Println("reading src file:", err.Error)
		os.Exit(1)
	} else {
		fmt.Println("the src file has been read")
	}

	err = ParseSrcData(srcData, &logs) // pass the read data and tmp storage for parsing
	if err != nil {
		fmt.Println("main(): Error:", err.Error())
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
