package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrParsingForm        = errors.New("error parsing form")
	ErrIncorrectSetDate   = errors.New("invalid start or end date")
	ErrIncorrectDateRange = errors.New("interval is more than 48 hours")
)

type fetchEntriesParams struct {
	// w                   http.ResponseWriter
	r                   *http.Request
	sourceNameParam     string    // key (key=value, e.g.: source_name=cute_ganymede)
	startDateParam      string    // key (key=value, e.g.: start_date=2022-01-01)
	endDateParam        string    // key (key=value, e.g.: end_date=2022-01-31)
	layoutDateTime      string    // e.g. "2006-01-02T15:04"
	startDate           time.Time // formatted start date
	endDate             time.Time // formatted end date
	rawQueryStr         string    // raw query string for make a SQL query
	dbName              string
	tableName           string
	timeStampColumnName string
	preparedQuery       string
}

func processURL(w http.ResponseWriter, r *http.Request) error {
	// checking whether the HTTP request method is a GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)

		return http.ErrNotSupported
	}

	// set headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // methods 'PUT', 'PATCH' and 'DELETE' has been deleted
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, hx-request, hx-target, hx-current-url")

	// parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return ErrParsingForm
	}

	return nil
}

func newFetchEntriesParams(w http.ResponseWriter, r *http.Request) *fetchEntriesParams {
	fep := new(fetchEntriesParams) // creating new fetchEntriesParams instance
	// fep.w = w
	fep.r = r

	// set default values
	fep.sourceNameParam = "source_name"     // default source name parameter (server's name)
	fep.startDateParam = "start_date"       // default start date and time parameter
	fep.endDateParam = "end_date"           // default end date and time parameter
	fep.layoutDateTime = "2006-01-02T15:04" // default layout for date and time

	fep.rawQueryStr = `SELECT * FROM %s WHERE %s BETWEEN '%s' AND '%s'` // building query string for SQL query
	fep.dbName = r.FormValue(fep.sourceNameParam)                       // getting source name URL parameters
	fep.tableName = "lg_tab"                                            // default table name in the database
	fep.timeStampColumnName = "tmstmp"                                  // default timestamp column name in the table

	return fep
}

func (f *fetchEntriesParams) parseAndValidateDateRange() error {

	startDate, err := time.Parse(f.layoutDateTime, f.r.FormValue(f.startDateParam))
	if err != nil {
		return ErrIncorrectSetDate
	}
	endDate, err := time.Parse(f.layoutDateTime, f.r.FormValue(f.endDateParam))
	if err != nil {
		return ErrIncorrectSetDate
	}

	// check if interval is valid (less than 48 hours)
	if endDate.Sub(startDate) > time.Hour*48 {
		return ErrIncorrectDateRange
	}

	f.startDate = startDate
	f.endDate = endDate

	return nil
}

func (f *fetchEntriesParams) buildFetchEntriesSQLString() {
	f.preparedQuery = fmt.Sprintf(
		f.rawQueryStr,
		f.tableName,
		f.timeStampColumnName,
		f.startDate.Format(f.layoutDateTime),
		f.endDate.Format(f.layoutDateTime),
	)
}
