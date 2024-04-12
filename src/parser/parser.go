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

	if log == "" {
		return LogEntry{}, ErrEmptyString
	}

	// split the input string into tokens, put them in a slice
	tokens := strings.Split(log, " ")

	for _, token := tokens {
		switch token {
			case
		}
	}
}
