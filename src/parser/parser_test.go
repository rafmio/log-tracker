package parser

import (
	"reflect"
	"testing"
	"time"
)

// Сначала создаем функцию, которая принимает строку файла - запись в журнале
// логов, парсим ее на структуру или карту

// type LogEntry struct {
// 	// Date   time.Time
// 	SrcIP  string
// 	Len    string
// 	Ttl    string
// 	Id     string
// 	Spt    string
// 	Dpt    string
// 	Window string
// }

type tstCase struct {
	inputStr string
	logEntry LogEntry
}

func TestParseLog(t *testing.T) {

	tmStmp1, _ := time.Parse("2006-01-02 15:04:05", "2024-04-11 23:03:12")
	tmStmp2, _ := time.Parse("2006-01-02 15:04:05", "2024-01-31 22:48:52")
	tmStmp3, _ := time.Parse("2006-01-02 15:04:05", "2024-03-18 12:42:12")
	tstCases := []tstCase{
		{
			inputStr: "Apr 11 23:03:12 localhost kernel: [23604836.487395] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=212.192.158.71 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=248 ID=54321 PROTO=TCP SPT=44496 DPT=2099 WINDOW=65535 RES=0x00 SYN URGP=0 ",
			logEntry: LogEntry{
				TmStmp: tmStmp1,
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
			inputStr: "Jan 31 22:48:52 localhost kernel: [23603977.313498] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=205.210.31.172 DST=194.58.102.129 LEN=44 TOS=0x00 PREC=0x00 TTL=249 ID=54321 PROTO=TCP SPT=55098 DPT=50067 WINDOW=65535 RES=0x00 SYN URGP=0 ",
			logEntry: LogEntry{
				TmStmp: tmStmp2,
				SrcIP:  "205.210.31.172",
				Len:    "44",
				Ttl:    "249",
				Id:     "54321",
				Spt:    "55098",
				Dpt:    "50067",
				Window: "65535",
			},
		},
		{
			inputStr: "Mar 18 12:42:12 localhost kernel: [21493976.759473] [UFW BLOCK] IN=eth0 OUT= MAC=52:54:00:7c:d8:0f:fe:54:00:7c:d8:0f:08:00 SRC=104.156.155.7 DST=194.58.102.129 LEN=40 TOS=0x00 PREC=0x00 TTL=240 ID=33841 PROTO=TCP SPT=41818 DPT=777 WINDOW=1024 RES=0x00 SYN URGP=0 ",
			logEntry: LogEntry{
				TmStmp: tmStmp3,
				SrcIP:  "104.156.155.7",
				Len:    "40",
				Ttl:    "240",
				Id:     "33841",
				Spt:    "41818",
				Dpt:    "777",
				Window: "1024",
			},
		},
	}

	for _, tstCase := range tstCases {
		t.Run("run ParseLog", func(t *testing.T) {
			got, err := ParseLog(tstCase.inputStr)

			if !reflect.DeepEqual(got, tstCase.logEntry) {
				t.Errorf("ParseLog() = %v, want %v", got, tstCase.logEntry)
			}

			if err != nil {
				t.Errorf("ParseLog() error = %v", err)
			}
		})
	}
}
