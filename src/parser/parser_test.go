package parser

import (
	"reflect"
	"testing"
)

// Сначала создаем функцию, которая принимает строку файла - запись в журнале
// логов, парсим ее на структуру или карту

type LogEntry struct {
	// Date   time.Time
	SrcIP  string
	Len    string
	Ttl    string
	Id     string
	Spt    string
	Dpt    string
	Window string
}

type tstCase struct {
	inputStr string
	logEntry LogEntry
}

func TestParseLog(t *testing.T) {
	tstCases := []tstCase{
		{
			inputStr: "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 ",
			logEntry: LogEntry{
				SrcIP:  "212.192.158.71",
				Len:    "44",
				Ttl:    "248",
				Id:     "54321",
				Spt:    "44496",
				Dpt:    "2099",
				Window: "65535",
			},
		},
		{
			inputStr: "Apr 11 22:48:52 localhost kernel: [23603977.313498] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=205.210.31.172 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=249 ID=54321 PROTO=TCP SPT=55098 DPT=50067 WINDOW=65535 RES=0x00 SYN URGP=0 ",
			logEntry: LogEntry{
				SrcIP:  "205.210.31.172",
				Len:    "44",
				Ttl:    "249",
				Id:     "54321",
				Spt:    "55098",
				Dpt:    "50067",
				Window: "65535",
			},
		},
	}

	for _, tstCase := range tstCases {
		t.Run("run ParseLog", func(t *testing.T) {
			got, _ := ParseLog(tstCase.inputStr)

			if !reflect.DeepEqual(got, tstCase.logEntry) {
				t.Errorf("ParseLog() = %v, want %v", got, tstCase.logEntry)
			}
		})
	}
}
