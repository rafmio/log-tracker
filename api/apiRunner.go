package main

import (
	"database/sql"
	"fmt"
	"log"
	"logtracker/dbhandler"
	"net/http"
)

const (
	dbConfigFileName = "db-config.json"
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

	// get parameters from form
	sourceName := r.FormValue("sourceName") // to know which of the two servers to send the request to
	// serviceType := r.FormValue("serviceType")
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")

	// read config file for connecting to servers
	sourceConfigs := readConfig(dbConfigFileName)
	sourceConfig := sourceConfigs[sourceName]

	// creating variables for storing ip and port
	var ip string
	var port string

	// read config file for connecting to DB
	dbConfigStruct, err := dbhandler.LoadDatabaseConfig(fileDBConfigName)
	if err != nil {
		log.Fatal(err)
	}

	dbConfigTxt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfigStruct.Host,
		dbConfigStruct.Port,
		dbConfigStruct.User,
		dbConfigStruct.Password,
		dbConfigStruct.Dbname,
		dbConfigStruct.Sslmode,
	)

	// connect to DB
	db, err := sql.Open(dbConfigStruct.DriverName, dbConfigTxt)
	if err != nil {
		log.Fatal(err)
	}

	// The date data is in the dd/mm/yyyy format, they need to be converted to the yyyy-mm-dd format
	startTime := startDate + " 00:00:00"
	endTime := endDate + " 23:59:59"

	// making a query to the databasequery
	rows, err := db.Query("SELECT * FROM lb_tab WHERE tmstmp BETWEEN $1 AND $1", startTime, endTime)
	if err != nil {
		log.Fatal(err)
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
