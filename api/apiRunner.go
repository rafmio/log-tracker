package main

import (
	"log"
	"net/http"
)

const (
	portTest                 string = "8081"
	portBlackOxygenium       string = "5432"
	portCuteGanymede         string = "5432"
	serverNameBlackOxygenium string = "BlackOxygenium"
	serverNameCuteGanymede   string = "CuteGanymede"
	ipBlackOxygenium         string = "194.58.102.129"
	ipCuteGanymede           string = "147.45.226.19"
	fileDBConfigName         string = "../config/databaseConfig.json"
	fileConfigName           string = "../config/fileConfig.json"
)

// fetchHandler() fetches data from PostgreSQL DB with given parameters
func fetchHandler(w http.ResponseWriter, r *http.Request) {

	// check if http method is 'GET'
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// parser form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// set headers
	w.Header().Set("Content-Type", "application/json")

	sourceName := r.FormValue("sourceName")
	serviceType := r.FormValue("serviceType")
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")

	// creating variables for storing ip and port
	var ip string
	var port string

	// setting the values of 'ip' and 'port' depending on the value of 'SourceName'
	switch sourceName {
	case serverNameBlackOxygenium:
		ip = ipBlackOxygenium
		port = portBlackOxygenium
	case serverNameCuteGanymede:
		ip = ipCuteGanymede
		port = portCuteGanymede
	default:
		http.Error(w, "Invalid sourceName", http.StatusBadRequest)
		return
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
