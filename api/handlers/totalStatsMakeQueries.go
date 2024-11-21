package handlers

import (
	"database/sql"
	"sync"
)

// func (s *serverStats) makeQueries() error {
// 	// initialize map for storing result of SQL queries (*sql.Rows)
// 	s.statIndicatorsRows = make(map[string]*sql.Rows, len(s.statIndicatorsQueries))

// 	// qrNm - key, the name of the query
// 	// sqlQr - value, the SQL query itself
// 	for qrNm, sqlQr := range s.statIndicatorsQueries {
// 		resultOfSqlQr, err := s.DB.Query(sqlQr)
// 		if err != nil {
// 			return err
// 		}
// 		s.statIndicatorsRows[qrNm] = resultOfSqlQr
// 	}

// 	return nil
// }

func (s *serverStats) makeQueries() error {
	// Initialize map for storing result of SQL queries (*sql.Rows)
	s.statIndicatorsRows = make(map[string]*sql.Rows, len(s.statIndicatorsQueries))

	// Create a channel to receive results and errors
	resultsChan := make(chan struct {
		name string
		rows *sql.Rows
		err  error
	}, len(s.statIndicatorsQueries))

	var wg sync.WaitGroup // WaitGroup для ожидания завершения всех горутин

	// Iterate over the queries and launch a goroutine for each
	for qrNm, sqlQr := range s.statIndicatorsQueries {
		wg.Add(1) // Increment the WaitGroup counter

		go func(name string, query string) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			resultOfSqlQr, err := s.DB.Query(query)
			if err != nil {
				resultsChan <- struct {
					name string
					rows *sql.Rows
					err  error
				}{name: name, err: err} // Send error to channel
				return
			}
			resultsChan <- struct {
				name string
				rows *sql.Rows
				err  error
			}{name: name, rows: resultOfSqlQr} // Send result to channel
		}(qrNm, sqlQr) // Pass parameters to the goroutine
	}

	go func() {
		wg.Wait()          // Wait for all goroutines to finish
		close(resultsChan) // Close the channel after all have finished
	}()

	// Collect results from the channel
	for result := range resultsChan {
		if result.err != nil {
			return result.err // Return the first error encountered
		}
		s.statIndicatorsRows[result.name] = result.rows // Store successful results
	}

	return nil
}

// func (s *serverStats) makeQueries() error {
// 	// Initialize map for storing result of SQL queries (*sql.Rows)
// 	s.statIndicatorsRows = make(map[string]*sql.Rows, len(s.statIndicatorsQueries))

// 	// Create a channel to receive results and errors
// 	resultsChan := make(chan struct {
// 		name string
// 		rows *sql.Rows
// 		err  error
// 	}, len(s.statIndicatorsQueries))

// 	var wg sync.WaitGroup // WaitGroup для ожидания завершения всех горутин

// 	// Iterate over the queries and prepare statements
// 	for qrNm, sqlQr := range s.statIndicatorsQueries {
// 		wg.Add(1) // Increment the WaitGroup counter

// 		go func(name string, query string) {
// 			defer wg.Done() // Decrement the counter when the goroutine completes

// 			stmt, err := s.DB.Prepare(query) // Prepare the statement
// 			if err != nil {
// 				resultsChan <- struct {
// 					name string
// 					rows *sql.Rows
// 					err  error
// 				}{name: name, err: err} // Send error to channel
// 				return
// 			}
// 			defer stmt.Close() // Ensure the statement is closed after use

// 			resultOfSqlQr, err := stmt.Query() // Execute the prepared statement
// 			if err != nil {
// 				resultsChan <- struct {
// 					name string
// 					rows *sql.Rows
// 					err  error
// 				}{name: name, err: err} // Send error to channel
// 				return
// 			}
// 			resultsChan <- struct {
// 				name string
// 				rows *sql.Rows
// 				err  error
// 			}{name: name, rows: resultOfSqlQr} // Send result to channel
// 		}(qrNm, sqlQr) // Pass parameters to the goroutine
// 	}

// 	go func() {
// 		wg.Wait()          // Wait for all goroutines to finish
// 		close(resultsChan) // Close the channel after all have finished
// 	}()

// 	// Collect results from the channel
// 	for result := range resultsChan {
// 		if result.err != nil {
// 			return result.err // Return the first error encountered
// 		}
// 		s.statIndicatorsRows[result.name] = result.rows // Store successful results
// 	}

// 	return nil
// }
