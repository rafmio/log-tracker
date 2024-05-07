package dbhandler

var filepath string = "../config/databaseConfig.json"

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
// var CDBc ConnectDBConfig = ConnectDBConfig{
// 	driverName: "postgres",
// 	user:       "raf",
// 	dbname:     "logtracker",
// 	tableName:  "lg_tab_1",
// 	password:   "qwq121",
// 	sslmode:    "disable",
// }

func LoadDatabaseConfig(filepath string) (ConnectDBConfig, error) {

}
