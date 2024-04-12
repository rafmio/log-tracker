package parser

import (
	"errors"
	"strings"
)

var ErrEmptyString = errors.New("log-string is empty")

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

func ParseLog(log string) (LogEntry, error) {

	result := LogEntry{}

	if log == "" {
		return LogEntry{}, ErrEmptyString
	}

	tokens := strings.Split(log, " ") // split log-entry into slices

	// DATE:
	// dateStr := tokens[:3] // first 3 elements of 'tokens are the date
	// tokens := tokens[3:] // trim 'tokens' - delete date

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
