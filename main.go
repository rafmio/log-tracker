package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var SrcFileName string = "ufw.log"

type Log struct {
	Date      time.Time
	HMtime    time.Time
	SrcIP     net.IP
	PacketLen int
	Ttl       int
	PacketId  int
	SrcPort   int
	DstPort   int
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

}
