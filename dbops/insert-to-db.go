package dbhandler

import (
	"database/sql"
	"fmt"
	"log"

	"logtracker/parser"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InsertToDb() connect to DB, check if the record exists, if not, insert the record
func InsertToDb(logEntry parser.LogEntry, cDBc ConnectDBConfig) error {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cDBc.Host,
		cDBc.Port,
		cDBc.User,
		cDBc.Dbname,
		cDBc.Password,
		cDBc.Sslmode,
	)

	db, err := sql.Open(cDBc.DriverName, dataSourceName)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer db.Close()

	doesRecordExists, err := CheckIfRecordExists(db, cDBc.TableName, logEntry)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if doesRecordExists {
		log.Println("the record already exists")
		return nil
	}

	// preparing the query for INSERT
	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		cDBc.TableName,
		"tmstmp",
		"srcip",
		"len",
		"ttl",
		"innerid",
		"spt",
		"dpt",
		"wndw",
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		logEntry.TmStmp,
		logEntry.SrcIP,
		logEntry.Len,
		logEntry.Ttl,
		logEntry.Id,
		logEntry.Spt,
		logEntry.Dpt,
		logEntry.Window,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// check if the record exists
func CheckIfRecordExists(db *sql.DB, tableName string, logEntry parser.LogEntry) (bool, error) {
	// preparing the query for SELECT
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s AND %s AND %s AND %s AND %s AND %s AND %s AND %s",
		tableName,
		"tmstmp = $1",
		"srcip = $2",
		"len = $3",
		"ttl = $4",
		"innerid = $5",
		"spt = $6",
		"dpt = $7",
		"wndw = $8",
	)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("preparing query: ", err.Error())
		return false, err
	}
	defer stmt.Close()

	// executing the query
	rows, err := stmt.Query(
		logEntry.TmStmp,
		logEntry.SrcIP,
		logEntry.Len,
		logEntry.Ttl,
		logEntry.Id,
		logEntry.Spt,
		logEntry.Dpt,
		logEntry.Window,
	)

	if err != nil {
		log.Println("executing query: ", err.Error())
		return false, err
	}
	defer rows.Close()

	// returns true on success, or false if there is no next result row
	// or an error happened while preparing it
	return rows.Next(), nil
}
