package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

type ltGeneralStats struct {
	dbConfigFilePath   string             // the path to the database connection configuration file
	dbConfigs          map[string]Source  // map database connection configurations returned by readConfig()
	dsn                string             // data source name
	db                 *sql.DB            // connection returned by sql.Open()
	tmpl               *template.Template // html template parsed from a file
	totalNumberEntries map[string]int
	uniqueIpCount      map[string]int
}

func (d *ltGeneralStats) setDbConfigFilePath() {
	d.dbConfigFilePath = "db-config.json"
}

// the final result of this functions - HTML code of three-column table
func generalStatHandler(w http.ResponseWriter, r *http.Request) {

	// DEBUG:
	w.Write([]byte("<p>DEBUG: hello from generalStatHandler</p>"))

	// read template with general statistic data
	tmpl, err := template.ParseFiles("generalStatTable.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // message to browser
		log.Println("Error parsing files:", err.Error())           // message to terminal (or logger)
	}
	tmpl = template.Must(tmpl, err)

	// dbConfigs := readConfig(dbConfigFileName)
	dbConfigs, err := readConfig(dbConfigFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // message to browser
		log.Println("Error reading config file:", err.Error())     // message to terminal (or logger)
	}

	// instances for data storing
	var data LtGeneralStats
	data.TotalNumberEntries = make(map[string]int)
	data.UniqueIpCount = make(map[string]int)

	// range list of configs to make query
	for _, dbConfig := range dbConfigs {
		db, err := openDB(dbConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // message to browser
			log.Println("Error opening database:", err.Error())        // message to terminal (or logger)
		}

		defer db.Close()

		data.TotalNumberEntries[dbConfig.Name] = getTotalNumberOfEntries(db)
		data.UniqueIpCount[dbConfig.Name] = getUniqueIpCount(db)
	}

	db, err := openDB(dbConfigs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // message to browser
		log.Println("Error opening database:", err.Error())        // message to terminal (or logger)
	}

	defer db.Close()

}
