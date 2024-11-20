package handlers

import (
	"database/sql"
)

func (s *serverStats) makeQueries() error {
	// initialize map for storing result of SQL queries (*sql.Rows)
	s.statIndicatorsRows = make(map[string]*sql.Rows, len(s.statIndicatorsQueries))

	// qrNm - key, the name of the query
	// sqlQr - value, the SQL query itself
	for qrNm, sqlQr := range s.statIndicatorsQueries {
		resultOfSqlQr, err := s.DB.Query(sqlQr)
		if err != nil {
			return err
		}
		s.statIndicatorsRows[qrNm] = resultOfSqlQr
	}

	return nil
}
