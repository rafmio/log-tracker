package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func openDB(srcConf *Source) (DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		srcConf.Host,
		srcConf.Port,
		srcConf.User,
		srcConf.DBName,
		srcConf.Password,
		srcConf.SslMode,
	))

	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to DB: %v", err), http.StatusInternalServerError)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Error pinging DB: %v", err), http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}
