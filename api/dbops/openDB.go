package dbops

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// import the PostgreSQL driver for database/sql
	_ "github.com/lib/pq" // $ go get .
)

type DBConfig struct {
	Name        string `json:"Name"`        // server's name for internal using in code, mapping etc ('cute_ganymede')
	DisplayName string `json:"DisplayName"` // the same name as 'Name', only for display ('Cute Ganymede')
	DriverName  string `json:"DriverName"`  // e.g. "postgres"
	Host        string `json:"Host"`        // "194.58.102.129", "localhost", etc
	Port        string `json:"Port"`        // port number, e.g. "5432", "8543", etc
	DBName      string `json:"DBName"`      // name of DB inside of 'PostgreSQL'
	User        string `json:"User"`        // username "raf", "postgres", etc
	Password    string `json:"Password"`    // password
	SslMode     string `json:"SslMode"`     // SSL mode, etc "disable", "require", "verify-full", etc"
	Dsn         string // data source name
	DB          *sql.DB
	Err         error
}

func NewDBConfig(dbConfigFilePath, srcName string) (*DBConfig, error) {

	// check if dbConfigFilePath is empty
	if dbConfigFilePath == "" {
		log.Println("Database config file path is empty")
		return nil, fmt.Errorf("database config file path is empty")
	}

	// reading file with configuration for DB connection
	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
		return nil, err
	}

	// unmarshalling JSON data to struct
	dbConfigs := make(map[string]DBConfig) // variable for storing unmarshalled data
	err = json.Unmarshal(file, &dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
		return nil, err
	}

	dbCfg := new(DBConfig)       // new DBConfig instance
	tmpCfg := dbConfigs[srcName] // assigning data from map to struct for given source name
	*dbCfg = tmpCfg              // assigning data from map to struct for given source name
	dbCfg.Name = srcName

	return dbCfg, nil
}

func (dbc *DBConfig) SetDSN() {
	formatString := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	dbc.Dsn = fmt.Sprintf(formatString,
		dbc.Host,
		dbc.Port,
		dbc.User,
		dbc.DBName,
		dbc.Password,
		dbc.SslMode,
	)
}

func (dbc *DBConfig) EstablishDbConnection() error {
	var err error
	dbc.DB, err = sql.Open(dbc.DriverName, dbc.Dsn)
	if err != nil {
		log.Println("Open database:", err)
		dbc.Err = err
		return err
	}

	err = dbc.DB.Ping()
	if err != nil {
		log.Println("Ping database:", err)
		dbc.Err = err
	}
	return nil
}
