package dbops

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	// import the PostgreSQL driver for datebase/sql
	_ "github.com/lib/pq" // $ go get .
)

// var (
// 	ErrDriverNameEmpty    = errors.New("driver name is empty")
// 	ErrDSNMapEmpty        = errors.New("DSN map is empty or nil")
// 	ErrOpeningDatabase    = errors.New("error opening database")
// 	ErrPingingDatabase    = errors.New("error pinging database")
// 	ErrOpenDBErrsMapEmpty = errors.New("openDBErrsMap is empty or nil")
// )

type DBConfig struct {
	Name             string `json:"Name"` // server's name for internal using in code, mapping etc ('cute_ganymede')
	DbConfigFilePath string
	DisplayName      string `json:"DisplayName"` // the same name as 'Name', only for display ('Cute Ganymede')
	DriverName       string `json:"DriverName"`  // e.g. "postgres"
	Host             string `json:"Host"`        // "194.58.102.129", "localhost", etc
	Port             string `json:"Port"`        // port number, e.g. "5432", "8543", etc
	DBName           string `json:"DBName"`      // name of DB inside of 'PostgreSQL'
	User             string `json:"User"`        // username "raf", "postgres", etc
	Password         string `json:"Password"`    // password
	SslMode          string `json:"SslMode"`     // SSL mode, etc "disable", "require", "verify-full", etc"
	Dsn              string // data source name
	DB               *sql.DB
	Err              error
}

func NewDBConfig(dbConfigFilePath string) (*DBConfig, error) {

	// check if dbConfigFilePath is empty
	if d.DbConfigFilePath == "" {
		log.Println("Database config file path is empty")
		return fmt.Errorf("Database config file path is empty")
	}

	// reading file with configuration for DB connection
	file, err := os.ReadFile(d.DbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
		return err
	}

	// unmarshalling JSON data to struct
	dbConfigs = make(map[string]DBConfig) // variable for storing unmarshalled data
	err = json.Unmarshal(file, &d.dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return err
	}

	dbCfg := new(DBConfig)
	dbCfg = &dbConfigs[d.Name]

	return dbCfg, nil
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
