package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	// import the PostgreSQL driver for datebase/sql
	_ "github.com/lib/pq" // $ go get .
)

type LogEntry struct {
	SeqNum string
	TmStmp time.Time
	SrcIP  string
	Len    string
	Ttl    string
	Id     string // will named 'inner id' in database
	Spt    string
	Dpt    string
	Window string // will named 'wndw' in database
}

const (
	dbConfigFileName = "db-config.json"
	port             = ":8082"
	sourceNameParam  = "source_name"
	startDateParam   = "start_date"
	endDateParam     = "end_date"
	layoutDateTime   = "2006-01-02T15:04"
)

/*
fetchHandler() handles incoming HTTP requests.

The format of the received request (example):
https://194.58.102.129:8082/fetch?source_name=cute_ganymede&start_date=2024-08-21T14:35&end_date=2024-08-22T11:50

parameter names:
- sourceName: Name of the source (black_oxygenium or cute_ganymede)
- start_date: Start date and time (ISO 8601) of the data to fetch
- end_date: End date and time (ISO 8601) of the data to fetch
*/
func fetchHandler(w http.ResponseWriter, r *http.Request) {

	// SET HEADERS ----------------------------------------------
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // methods 'PUT', 'PATCH' and 'DELETE' has been deleted
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, hx-request, hx-target, hx-current-url")

	// parse form for determine source (exact DB) and further date parsing
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// determine which source (server) received the request
	// get parameters from Request.Form
	sourceNamesStr := r.FormValue(sourceNameParam)       // to know which of the two servers to send the request to
	sourceNamesSls := strings.Split(sourceNamesStr, ",") // sourceNamesSls is of type []string

	// slice for storing list of configs of DB connections to servers (sources)
	// dbConfigs := readConfig(dbConfigFileName) // returns map[string]Source
	dbConfigs, err := readConfig(dbConfigFileName) // returns map[string]Source
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config file: %v", err), http.StatusInternalServerError)
		return
	}

	for _, src := range sourceNamesSls {

		// if the source exists in the configuration
		if _, ok := dbConfigs[src]; !ok {
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}

		// set connection to DB
		srcConf := dbConfigs[src] // 'srcConf' if of type Source

		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			srcConf.Host,
			srcConf.Port,
			srcConf.User,
			srcConf.DBName,
			srcConf.Password,
			srcConf.SslMode,
		))

		if err != nil {
			http.Error(w, fmt.Sprintf("Error connecting to DB: %v", err), http.StatusInternalServerError)
			return
		}

		if err = db.Ping(); err != nil {
			http.Error(w, fmt.Sprintf("Error pinging DB: %v", err), http.StatusInternalServerError)
			return
		}

		defer db.Close() // close connection to DB

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

		// query DB for data between startDate and endDate
		tmstmpColumnName := "tmstmp"
		query := fmt.Sprintf(`SELECT * FROM %s WHERE %s >= $1 AND %s <= $2`,
			srcConf.TableName,
			tmstmpColumnName,
			tmstmpColumnName)

		rows, err := db.Query(
			query,
			startDate, // $3
			endDate,   // $4
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying DB: %v", err), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		entries := make([]LogEntry, 0)

		for rows.Next() {
			var entry LogEntry
			err := rows.Scan(
				&entry.SeqNum,
				&entry.TmStmp,
				&entry.SrcIP,
				&entry.Len,
				&entry.Ttl,
				&entry.Id,
				&entry.Spt,
				&entry.Dpt,
				&entry.Window,
			)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}

			entries = append(entries, entry) // entries is of type []LogEntry
		}
		w.WriteHeader(http.StatusOK)

		tmpl, err := template.New("logTable").Parse(htmlTemplateLogsTable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		// var htmlOutput string
		err = tmpl.Execute(w, entries)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println("DEBUG: end of fetchHandler(). Response headers:")
		fmt.Println(w.Header())
		// fmt.Fprintf(w, htmlOutput)
	}
}

func main() {
	mux := http.NewServeMux()

	fmt.Printf("Listening on port %s\n", port)

	// creating routes
	mux.HandleFunc("/fetch", fetchHandler)
	mux.HandleFunc("/logtracker/statistic/generalstat", generalStatHandler)
	// running server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
