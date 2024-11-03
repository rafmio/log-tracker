package main

import (
	"net/http"
	"time"
)

type LogEntry struct {
	SeqNum string
	TmStmp time.Time // timestamp
	SrcIP  string    // source IP address
	Len    string
	Ttl    string
	Id     string // will named 'innerid' in database
	Spt    string // source port
	Dpt    string // destination port
	Window string // will named 'wndw' in database
}

/*
fetchEntriesHandler() handles incoming HTTP requests.

The format of the received request (example):
https://194.58.102.129:8082/logtracker/fetch_entries?source_name=cute_ganymede&start_date=2024-08-21T14:35&end_date=2024-08-22T11:50

parameter names:
- source_name: Name of the source (black_oxygenium or cute_ganymede)
- start_date: Start date and time (ISO 8601) of the data to fetch
- end_date: End date and time (ISO 8601) of the data to fetch
*/
func fetchEntriesHandler(w http.ResponseWriter, r *http.Request) {
	err := processUrl(w, r)
	if err != nil {
		if err == http.ErrNotSupported {
			http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)
		}
		if err == ErrParsingForm {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}
	}

<<<<<<< HEAD
	fep := newFetchEntriesParams(w, r)

	err = fep.parseAndValidateDateRange()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
=======
>>>>>>> 94b458019f19ce295a241612a698ab8922e96574
}
