package main

import (
	"fmt"
	"net/http"
	"time"
)

type LogEntry struct {
	SeqNum string
	TmStmp time.Time
	SrcIP  string
	Len    string
	Ttl    string
	Id     string // will named 'innerid' in database
	Spt    string
	Dpt    string
	Window string // will named 'wndw' in database
}

const (
	sourceNameParam = "source_name"
	startDateParam  = "start_date"
	endDateParam    = "end_date"
	layoutDateTime  = "2006-01-02T15:04"
)

type queryParams struct {
	sourceNameParam string
	startDateParam  string
	endDateParam    string
	layoutDateTime  string
	startDate       time.Time
	endDate         time.Time
}

/*
fetchEntriesHandler() handles incoming HTTP requests.

The format of the received request (example):
https://194.58.102.129:8082/logtracker/fetch_entries?source_name=cute_ganymede&start_date=2024-08-21T14:35&end_date=2024-08-22T11:50

parameter names:
- source_ame: Name of the source (black_oxygenium or cute_ganymede)
- start_date: Start date and time (ISO 8601) of the data to fetch
- end_date: End date and time (ISO 8601) of the data to fetch
*/
func fetchEntriesHandler(w http.ResponseWriter, r *http.Request) {
	// preProcessResponse make:
	// 	1. checking whether the HTTP request method is a GET method
	// 	2. setting headers
	// 	3. parsing request's parameters and populate r.Form
	err := preProcessResponse(w, r)
	if err != nil {
		if err == http.ErrNotSupported {
			http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)
		}
		if err == ErrParsingForm {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}
	}

	startDate, err := time.Parse(layoutDateTime, r.FormValue(startDateParam))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid start_date: %v", err), http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse(layoutDateTime, r.FormValue(endDateParam))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid end_date: %v", err), http.StatusBadRequest)
		return
	}

	// check if interval is valid (less than 48 hours)
	if endDate.Sub(startDate) > time.Hour*48 {
		http.Error(w, "<p>Interval is too long.Please set the interval less than 48 hours</p>", http.StatusBadRequest)
		return
	}
}
