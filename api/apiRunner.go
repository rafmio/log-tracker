package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// import the PostgreSQL driver for datebase/sql
	_ "github.com/lib/pq" // $ go get .
)

const (
	dbConfigFileName = "db-config.json"
	port             = ":8082"

	// variables URL parameter names
	// sourceName = "sourceName"
	// startDate = "startDate"
	// endDate = "endDate"
)

type DateTimeRange struct {
	StartDate time.Time
	EndDate   time.Time
	StartTime time.Time
	EndTime   time.Time
}

/*
fetchHandler() handles incoming HTTP requests.

It performs the following tasks:
- checks if the HTTP method is 'GET'. If not, it responds with a 405 status code and an error message.
- parses the form data from the request
- set headers
- determines which source (server) received the request: put list of servers in the slice
*/
func fetchHandler(w http.ResponseWriter, r *http.Request) {

	// check if http method is 'GET'
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// set headers
	w.Header().Set("Content-Type", "application/json")

	// determine which source (server) received the request
	// get parameters from Request.Form
	sourceNamesStr := r.FormValue("sourceName")          // to know which of the two servers to send the request to
	sourceNamesSls := strings.Split(sourceNamesStr, ",") // sourceNamesSls is of type []string

	// slice for storing list of configs of DB connections to servers (sources)
	dbConfigs := readConfig(dbConfigFileName) // returns map[string]Source

	for _, src := range sourceNamesSls {
		// log.Println("src:", src) // for debugging
		if _, ok := dbConfigs[src]; !ok {
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}

		// for sector debugging:
		// srvr := dbConfigs[src]
		// w.Write([]byte(srvr.Name))
		// w.Write([]byte("\n"))
		// w.Write([]byte(srvr.Host))
		// w.Write([]byte("\n"))
		// end of debugging sector

		// set connection to DB
		srcConf := dbConfigs[src] // srcConf if of type Source

		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			srcConf.Host,
			srcConf.Port,
			srcConf.User,
			srcConf.DBName,
			srcConf.Password,
			srcConf.SslMode,
		))

		if err != nil {
			// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("Error connecting to DB: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		startDateStr := r.FormValue("startDate")
		endDateStr := r.FormValue("endDate")

		// parse dates
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid startDate format", http.StatusBadRequest)
			return
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid endDate format", http.StatusBadRequest)
			return
		}

		queryString := "SELECT * FROM " + srcConf.TableName + " WHERE tmstmp BETWEEN $1 AND $2"

		rows, err := db.Query(queryString, startDate, endDate)
		if err != nil {
			http.Error(w, "Error executing query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var columnValue string
			err = rows.Scan(&columnValue)
			if err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(columnValue + "\n"))
		}
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
