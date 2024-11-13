package handlers

import (
	"api/dbops"
	"fmt"
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
func FetchEntriesHandler(w http.ResponseWriter, r *http.Request) {
	// preProcessResponse make:
	// 	1. checking whether the HTTP request method is a GET method
	// 	2. setting headers
	// 	3. parsing request's parameters and populate r.Form
	err := processURL(w, r)
	if err != nil {
		if err == http.ErrNotSupported {
			http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)
		}
		if err == ErrParsingForm {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}
	}

	// new instance
	fep := newFetchEntriesParams(w, r)

	err = fep.parseAndValidateDateRange()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// build SQL query
	fep.buildFetchEntriesSQLString()

	// DEBUG print
	fmt.Fprintf(w, "<p>start_date: %s</p>", fep.startDate)
	fmt.Fprintf(w, "<p>end_date: %s</p>", fep.endDate)
	fmt.Fprintf(w, "<p>prepared query: %s</p>", fep.preparedQuery)
	// end of DEBUG

	dbCfg, err := dbops.NewDBConfig("config/db-config.json", fep.dbName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "<p>DEBUG: openDB.go</p>")
	fmt.Fprintf(w, "%s, %s, %s\n", dbCfg.DisplayName, dbCfg.Host, dbCfg.Port)
	fmt.Printf("%s, %s, %s\n", dbCfg.DisplayName, dbCfg.Host, dbCfg.Port)
	dbCfg.SetDSN()
	fmt.Fprintf(w, "<p>DSN: %s</p>", dbCfg.Dsn)

	// opening database connection
	err = dbCfg.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbCfg.DB.Close()

	rows, err := dbCfg.DB.Query(fep.preparedQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// iterate over rows and populate LogEntry structs
}
