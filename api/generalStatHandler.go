package main

import (
	"log"
	"net/http"
	"text/template"
)

// the final result of this functions - HTML code of three-column table
func generalStatHandler(w http.ResponseWriter, r *http.Request) {
	// read template with general statistic data
	tmpl, err := template.ParseFiles("generalStatTable.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // message to browser
		log.Println("Error parsing files:", err.Error())           // message to terminal (or logger)
	}
	tmpl = template.Must(tmpl, err)

	dbConfigs := readConfig(dbConfigFileName)

	// instances for data storing
	type Data struct {
		TotalNumberEntries map[string]int
		UniqueIpCount      map[string]int
	}

}
