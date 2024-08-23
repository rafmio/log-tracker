package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	// import the PostgreSQL driver for datebase/sql
	_ "github.com/lib/pq" // $ go get .
)

const (
	dbConfigFileName = "db-config.json"
	port             = ":8082"
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

	// check if http method is 'GET'
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse form for determine source (exact DB) and further date parsing
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// determine which source (server) received the request
	// get parameters from Request.Form
	sourceNamesStr := r.FormValue("source_name")         // to know which of the two servers to send the request to
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

		startDateStr := r.FormValue("start_date")
		endDateStr := r.FormValue("end_date")

	}

}

func main() {
	mux := http.NewServeMux()

	// creating routes
	mux.HandleFunc("/fetch", fetchHandler)

	// running server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
