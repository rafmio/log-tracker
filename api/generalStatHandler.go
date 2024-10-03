package main

import (
	"database/sql"
	"text/template"
)

type ltGeneralStats struct {
	dbConfigFilePath   string             // the path to the database connection configuration file - setDbConfigFilePath()
	dbConfigs          map[string]Source  // map database connection configurations returned by readConfig()
	dsns               map[string]string  // data source names for every database, returned by setDSNs(). Format: map["server_name"]"formatString"
	dbs                map[string]*sql.DB // connections returned by sql.Open()
	tmpl               *template.Template // html template parsed from a file
	totalNumberEntries map[string]int
	uniqueIpCount      map[string]int
}
