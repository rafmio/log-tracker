package dbhandler

import (
	"encoding/json"
	"log"
	"os"
)

var ConfDBFilePath string = "/home/raf/log-tracker/config/databaseConfig.json"

// declare the structure of the database connection parameters
type ConnectDBConfig struct {
	DriverName string
	User       string
	Dbname     string
	TableName  string
	Password   string
	Sslmode    string
}

// initialize ConnectDBConfig with exact values
// var CDBc ConnectDBConfig = ConnectDBConfig{
// 	driverName: "postgres",
// 	user:       "raf",
// 	dbname:     "logtracker",
// 	tableName:  "lg_tab_1",
// 	password:   "qwq121",
// 	sslmode:    "disable",
// }

// Extracts the settings from the json configuration file and returns
// the ConnectDBConfig structure
func LoadDatabaseConfig(ConfDBFilePath string) (ConnectDBConfig, error) {
	data, err := os.ReadFile(ConfDBFilePath)
	if err != nil {
		log.Println("reading DB configs:", err)
		return ConnectDBConfig{}, err
	}

	var CDBc ConnectDBConfig

	err = json.Unmarshal(data, &CDBc)
	if err != nil {
		return ConnectDBConfig{}, err
	}

	return CDBc, nil
}
