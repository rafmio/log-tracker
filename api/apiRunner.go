package main

import (
	"log"
	"net/http"
)

const (
	dbConfigFileName = "db-config.json"
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
	sourceName := r.FormValue("sourceName") // to know which of the two servers to send the request to

	// slice for storing list of configs of DB connections to servers (sources)
	dbConfigs := readConfig(dbConfigFileName)

	for key, val := range dbConfigs {
		w.Write([]byte(val))
	}
}

func main() {
	mux := http.NewServeMux()

	// creating routes
	mux.HandleFunc("/fetch", fetchHandler)

	// running server
	if err := http.ListenAndServe(portTest, mux); err != nil {
		log.Fatal(err)
	}
}
