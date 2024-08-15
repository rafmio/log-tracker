package main

import (
	"log"
	"net/http"
	"strings"
)

const (
	dbConfigFileName = "db-config.json"
	port             = ":8082"

	// variables URL parameter names
	// sourceName = "sourceName"
	// startDate = "startDate"
	// endDate = "endDate"
)

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
	sourceNamesStr := r.FormValue("sourceName") // to know which of the two servers to send the request to
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
		startDateStr := r.FormValue()

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
