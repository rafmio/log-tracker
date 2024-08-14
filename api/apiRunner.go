package main

import (
	"log"
	"net/http"
	"strings"
)

const (
	dbConfigFileName = "db-config.json"
	port             = ":8082"
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
	sourceNamesSls := strings.Split(sourceNamesStr, ",")

	// slice for storing list of configs of DB connections to servers (sources)
	dbConfigs := readConfig(dbConfigFileName) // returns map[string]Source

	for _, src := range sourceNamesSls {

		if _, ok := dbConfigs[src]; !ok {
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}

		srvr := dbConfigs[src]
		w.Write([]byte(srvr.Name))
		w.Write([]byte("\n"))
		w.Write([]byte(srvr.Host))
		w.Write([]byte("\n"))
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
