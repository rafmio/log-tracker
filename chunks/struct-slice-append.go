package main

import (
	"fmt"
	"os"
	"time"
)

type Log struct {
	Date   time.Time
	HMtime time.Time
	// SrcIP     net.IP
	// PacketLen int
	// Ttl       int
	// PacketId  int
	// SrcPort   int
	// DstPort   int
	// Window    int
}

func main() {
	slsInt := make([]int, 0)
	slsInt = append(slsInt, 150)
	ptrSlsInt := &slsInt
	slsAppend(ptrSlsInt)

	fmt.Println(slsInt)
}

func slsAppend(slsInt *[]int) {
	number := 100
	*slsInt = append(*slsInt, number)
	fmt.Fprintf(os.Stdout, "slsAppend(): %d\n", slsInt)
}
