package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	ErrEmptyString    = errors.New("log-string is empty")
	ErrParseTimeStamp = errors.New("can't parse time to timestamp")
)

type LogEntry struct {
	TmStmp time.Time
	SrcIP  string
	Len    string
	Ttl    string
	Id     string // will named 'inner id' in database
	Spt    string
	Dpt    string
	Window string // will named 'wndw' in database
}

func ParseLog(logLine string) (LogEntry, error) {

	result := LogEntry{}

	if logLine == "" {
		return LogEntry{}, ErrEmptyString
	}

	// tokens := strings.Split(logLine, " ") // split log-entry into slices
	tokens := strings.Fields(logLine)

	// PARSE TIMESTAMP
	// the UFW log file does not specify the year, we assign the current one
	year := time.Now().Year()
	yearStr := fmt.Sprint(year)

	months := make(map[string]string, 12)
	months["Jan"] = "01"
	months["Feb"] = "02"
	months["Mar"] = "03"
	months["Apr"] = "04"
	months["May"] = "05"
	months["Jun"] = "06"
	months["Jul"] = "07"
	months["Aug"] = "08"
	months["Sep"] = "09"
	months["Oct"] = "10"
	months["Nov"] = "11"
	months["Dec"] = "12"

	month := months[tokens[0]]

	// if the number of the month consists of one character, then add '0' before it
	if len(tokens[1]) == 1 {
		tokens[1] = "0" + tokens[1]
	}

	// cast to the appropriate format so that the time.Parse was able to parse
	timeStampStr := yearStr + "-" + month + "-" + tokens[1] + " " + tokens[2]

	timeStamp, err := time.Parse("2006-01-02 15:04:05", timeStampStr)
	if err != nil {
		log.Println("can't parse time") // TODO: make log system
		return LogEntry{}, ErrParseTimeStamp
	}
	result.TmStmp = timeStamp

	// PARSE LOG
	for _, token := range tokens {
		if strings.HasPrefix(token, "SRC=") {
			result.SrcIP = strings.TrimPrefix(token, "SRC=")
		} else if strings.HasPrefix(token, "LEN=") {
			result.Len = strings.TrimPrefix(token, "LEN=")
		} else if strings.HasPrefix(token, "TTL=") {
			result.Ttl = strings.TrimPrefix(token, "TTL=")
		} else if strings.HasPrefix(token, "ID=") {
			result.Id = strings.TrimPrefix(token, "ID=")
		} else if strings.HasPrefix(token, "SPT=") {
			result.Spt = strings.TrimPrefix(token, "SPT=")
		} else if strings.HasPrefix(token, "DPT=") {
			result.Dpt = strings.TrimPrefix(token, "DPT=")
		} else if strings.HasPrefix(token, "WINDOW=") {
			result.Window = strings.TrimPrefix(token, "WINDOW=")
		}
	}

	return result, nil
}
