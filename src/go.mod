module logtracker

go 1.22.2

replace logtracker/dbhandler => ./dbhandler

replace logtracker/parser => ./parser

require (
	logtracker/dbhandler v0.0.0-00010101000000-000000000000
	logtracker/parser v0.0.0-00010101000000-000000000000
)
