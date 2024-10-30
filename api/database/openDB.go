package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	// import the PostgreSQL driver for datebase/sql
	_ "github.com/lib/pq" // $ go get .
)

var (
	ErrDriverNameEmpty    = errors.New("driver name is empty")
	ErrDSNMapEmpty        = errors.New("DSN map is empty or nil")
	ErrOpeningDatabase    = errors.New("error opening database")
	ErrPingingDatabase    = errors.New("error pinging database")
	ErrOpenDBErrsMapEmpty = errors.New("openDBErrsMap is empty or nil")
)

type DBConnections struct {
	dbConfigFilePath string                     // the path to the database connection configuration file - setDbConfigFilePath()
	dbConfigs        map[string]ConnectDBConfig // map database connection configurations returned by readConfig()
	dsns             map[string]string          // data source names for every database, returned by setDSNs(). Format: map["server_name"]"formatString"
	dbs              map[string]*sql.DB         // connections returned by sql.Open() - don't forget to db.Close()!
	openDbErrs       map[string]error           // errors returned by sql.Open()
}

type ConnectDBConfig struct {
	DriverName  string `json:"DriverName"`  // e.g. "postgres"
	Name        string `json:"Name"`        // server's name for internal using in code, mapping etc ('cute_ganymede')
	DisplayName string `json:"DisplayName"` // the same name as 'Name', only for display ('Cute Ganymede')
	Host        string `json:"Host"`        // "194.58.102.129", "localhost", etc
	Port        string `json:"Port"`        // port number, e.g. "5432", "8543", etc
	DBName      string `json:"DBName"`      // name of DB inside of 'PostgreSQL'
	User        string `json:"User"`        // username "raf", "postgres", etc
	TableName   string `json:"TableName"`   // name of the table inside certain DB
	Password    string `json:"Password"`    // password
	SslMode     string `json:"SslMode"`     // SSL mode, etc "disable", "require", "verify-full", etc"
}

// Setting the path from where we will read the configuration file to connect to the database
func (dbC *DBConnections) setDBconfigFilePath() {
	dbC.dbConfigFilePath = "db-config.json"
}

func (dbC *DBConnections) readConfig() error {

	// check if dbConfigFilePath is empty
	if dbC.dbConfigFilePath == "" {
		log.Println("Database config file path is empty")
		return fmt.Errorf("Database config file path is empty")
	}

	// reading file with configuration for DB connection
	file, err := os.ReadFile(dbC.dbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
		return err
	}

	// unmarshalling JSON data to struct
	dbC.dbConfigs = make(map[string]ConnectDBConfig) // variable for storing unmarshalled data
	err = json.Unmarshal(file, &dbC.dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return err
	}

	return nil
}

func (dbC *DBConnections) setDSNs() {
	formatString := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	dbC.dsns = make(map[string]string)

	for _, dbConfig := range dbC.dbConfigs {
		dbC.dsns[dbConfig.Name] = fmt.Sprintf(formatString,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.DBName,
			dbConfig.Password,
			dbConfig.SslMode,
		)
	}
}

func (dbC *DBConnections) openDBs() error {

	// check if dsns is empty or nil
	if len(dbC.dsns) == 0 || dbC.dsns == nil {
		log.Println(ErrDSNMapEmpty)
		return ErrDSNMapEmpty
	}

	dbC.dbs = make(map[string]*sql.DB)      // variable for storing collection of DBs
	dbC.openDbErrs = make(map[string]error) // map[serverName]error

	var wg sync.WaitGroup
	var mu sync.Mutex

	// sql.Open() every source (database on certain server) in separate goroutine
	for serverName, dsn := range dbC.dsns {
		wg.Add(1)

		go func(serverName, dsn string) {
			driverName := dbC.dbConfigs[serverName].DriverName

			// check if dsns is empty or nil
			if driverName == "" {
				log.Println(ErrDriverNameEmpty)

				mu.Lock()
				dbC.openDbErrs[serverName] = ErrDriverNameEmpty
				mu.Unlock()
			}

			// open database
			db, err := sql.Open(driverName, dsn)
			if err != nil {
				log.Println(ErrOpeningDatabase)

				mu.Lock()
				dbC.openDbErrs[serverName] = err
				mu.Unlock()
			}

			// verify a connection to the database is still alive
			if err = db.Ping(); err != nil {
				log.Println(ErrPingingDatabase)

				mu.Lock()
				dbC.openDbErrs[serverName] = err
				mu.Unlock()
			}

			// add connection to collection of DBs
			mu.Lock()
			dbC.dbs[serverName] = db
			mu.Unlock()

			defer wg.Done()
		}(serverName, dsn)
	}

	wg.Wait()

	return nil // don't forget to db.Close()!
}
