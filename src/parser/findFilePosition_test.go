package parser

import (
	"testing"
)

// the functions under test gets *os.File as an argument,
// evaluates current file position, and save it in env variable.
// This is the test function for it:

func TestFindFilePosition(t *testing.T) {
	// test line for write into temp files
	line := "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 "

	// creating temp directory for temp files for run tests

}
