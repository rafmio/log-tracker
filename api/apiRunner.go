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

the format of the received request (example):
https://194.58.102.129:5432/fetch?source_name=cute_ganymede&start_date=2024-08-21T14:35&end_date=2024-08-22T11:50

parameter names:
- sourceName: Name of the source (black_oxygenium of cute_ganymede)
- start_date: Start date and time (ISO 8601) of the data to fetch
- end_date: End date and time (ISO 8601) of the data to fetch
*/
func fetchHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		fmt.Println("r.Method:", r.Method)
		// return // finish handling if it is preflight query
		// SET HEADERS ----------------------------------------------
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Установка заголовка CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешить все источники
		// Или, чтобы разрешить только конкретный источник:
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusOK)
		// ----------------------------------------------------------
	}

	// parse form for determine source (exact DB) and further date parsing
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	log.Println("Request received:", r.Form) // DEBUG INFO

	// determine which source (server) received the request
	// get parameters from Request.Form
	sourceNamesStr := r.FormValue(sourceNameParam)       // to know which of the two servers to send the request to
	sourceNamesSls := strings.Split(sourceNamesStr, ",") // sourceNamesSls is of type []string

	// slice for storing list of configs of DB connections to servers (sources)
	dbConfigs := readConfig(dbConfigFileName) // returns map[string]Source

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

		// close connection to DB after use
		defer db.Close()

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

		// query DB for data between startDate and endDate
		tmstmpColumnName := "tmstmp"
		// query := `SELECT * FROM $1 WHERE $2 >= $3 AND $2 <= $4` // move the line above to the block with constants
		query := fmt.Sprintf(`SELECT * FROM %s WHERE %s >= $1 AND %s <= $2`,
			srcConf.TableName,
			tmstmpColumnName,
			tmstmpColumnName)

		rows, err := db.Query(
			query,
			// srcConf.TableName, // $1
			// tmstmpColumnName,  // $2
			startDate, // $3
			// tmstmpColumnName,  // $2
			endDate, // $4
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

	// running server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
