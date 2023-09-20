package main

import (
	"net"
	"time"
)

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
