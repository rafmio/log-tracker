module logtracker

go 1.22.2

replace logtracker/dbhanlder => ./dbhanlder

replace logtracker/parser => ./parser

require logtracker/parser v0.0.0-00010101000000-000000000000 // indirect
