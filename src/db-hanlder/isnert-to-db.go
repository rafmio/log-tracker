package dbhanlder

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// declare the structure of the database connection parameters
type ConnectDBConfig struct {
	driverName string
	user       string
	dbname     string
	password   string
	sslmode    string
}

// initialize ConnectDBConfig with exact values
var cDBc ConnectDBConfig = ConnectDBConfig{
	driverName: "postgres",
	user:       "raf",
	dbname:     "logtracker",
	password:   "qwq121",
	sslmode:    "disable",
}

// InsertToDb() connect to DB, check if the record exists, if not, insert the record
func InsertToDb(logEntry *LogEntry, cDBc ConnectDBConfig) error {
	dataSourceName := fmt.Sprint("user=%s dbname=%s password=%s sslmode %s",
		cDBc.user,
		cDBc.dbname,
		cDBc.password,
		cDBc.sslmode,
	)

	db, err := sql.Open(cDBc.driverName, dataSourceName)
	if err != nil {
		log.Println(err.Error())
		return err
	} else {
		log.Println("db has been opened")
	}
	defer db.Close()

}
