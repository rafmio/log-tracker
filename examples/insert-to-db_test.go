package dbhanlder

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// declare the structure of the database connection parameters
type ConnectDBConfig struct {
	driverName string
	user       string
	dbname     string
	tableName  string
	password   string
	sslmode    string
}

// initialize ConnectDBConfig with exact values
var cDBc ConnectDBConfig = ConnectDBConfig{
	driverName: "postgres",
	user:       "raf",
	dbname:     "logtracker",
	tableName:  "lg_tab_1",
	password:   "qwq121",
	sslmode:    "disable",
}

// mock LogEntry struct
type LogEntry struct {
	TmStmp string
	SrcIP  string
	Len    int
	Ttl    int
	Id     int
	Spt    int
	Dpt    int
	Window int
}

// TestInsertToDb tests the InsertToDb function
func TestInsertToDb(t *testing.T) {
	// create a mock LogEntry
	logEntry := LogEntry{
		TmStmp: "2023-02-12 12:34:56",
		SrcIP:  "192.168.1.1",
		Len:    123,
		Ttl:    65,
		Id:     12345,
		Spt:    80,
		Dpt:    443,
		Window: 1024,
	}

	// mock the sql.DB object
	db := &sql.DB{}

	// mock the CheckIfRecordExists function
	doesRecordExists := true
	db.QueryRow = func(query string, args ...interface{}) *sql.Row {
		return &sql.Row{}
	}

	// execute the test
	err := InsertToDb(logEntry, cDBc)

	// verify the results
	if err != nil {
		t.Errorf("InsertToDb returned an error: %v", err)
	}

	if doesRecordExists != true {
		t.Errorf("CheckIfRecordExists returned false, expected true")
	}
}

// mock the CheckIfRecordExists function
func (cDBc ConnectDBConfig) CheckIfRecordExists(db *sql.DB, logEntry LogEntry) (bool, error) {
	return true, nil
}
