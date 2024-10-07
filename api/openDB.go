package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

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
func setDBconfigFilePath() string {
	return "db-config.json"
}

func readConfig(dbConfigFilePath string) (map[string]ConnectDBConfig, error) {

	// check if dbConfigFilePath is empty
	if dbConfigFilePath == "" {
		log.Println("The variable of the path to the file with the configuration of the database connection is empty")
		return nil, fmt.Errorf("The variable of the path to the file with the configuration of the database connection is empty")
	}

	// reading file with configuration for DB connection
	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
		return nil, err
	}

	// unmarshalling JSON data to struct
	dbConfigs := make(map[string]ConnectDBConfig) // variable for storing unmarshalled data
	err = json.Unmarshal(file, &dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return nil, err
	}

	return dbConfigs, err
}

func setDSNs(dsns map[string]string, dbConfigs map[string]ConnectDBConfig) {
	formatString := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	for _, dbConfig := range dbConfigs {
		dsns[dbConfig.Name] = fmt.Sprintf(formatString,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.DBName,
			dbConfig.Password,
			dbConfig.SslMode,
		)
	}
}

func openDBs(dsns map[string]string) (map[string]*sql.DB, map[string]error) {

	// check if driverName is empty
	if driverName == "" {
		log.Println("The variable of the name of the driver is empty")
		return nil, fmt.Errorf("The variable of the name of the driver is empty")
	}

	// check if dsns is empty or nil
	if len(dsns) == 0 || dsns == nil {
		log.Println("The variable of the map with the DSNs is empty or nil")
		return nil, fmt.Errorf("The variable of the map with the DSNs is empty or nil")
	}

	dbs := make(map[string]*sql.DB)      // variable for storing collection of DBs
	openDbErrs := make(map[string]error) // map[serverName]error

	var wg sync.WaitGroup
	var mu sync.Mutex

	// sql.Open() every source (database on certain server) in separate goroutine
	for serverName, dsn := range dsns {
		wg.Add(1)

		go func(serverName, dsn string) {
			db, err := sql.Open(driverName, dsn)
			if err != nil {
				log.Printf("Opening %s DB: %v\n", serverName, err)

				mu.Lock()
				openDbErrs[serverName] = err
				mu.Unlock()
			}
			// verify a connection to the database is still alive
			if err = db.Ping(); err != nil {
				log.Printf("Pinging %s DB: %v\n", serverName, err)

				mu.Lock()
				openDbErrs[serverName] = err
				mu.Unlock()
			}

			mu.Lock()
			dbs[serverName] = db
			mu.Unlock()

			defer wg.Done()
		}(serverName, dsn)
	}

	wg.Wait()

	if len(openDbErrs) == len(dsns) {
		return nil, fmt.Errorf("No database is available")
	} else if len(openDbErrs) > 0 && len(openDbErrs) < len(dsns) {
		return dbs, fmt.Errorf("Some databases are not available")
	}

	return dbs, nil // don't forget to db.Close()!
}
