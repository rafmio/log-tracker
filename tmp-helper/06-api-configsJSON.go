package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	// Dsn         string // data source name
	// DB          *sql.DB
	// Err         error
}

func main() {
	servers := make(map[string]DBConfig, 3)

	servers["black_oxygenium"] = DBConfig{
		Name:        "black_oxygenium",
		DisplayName: "Black Oxygenium",
		DriverName:  "postgres",
		Host:        "194.58.102.129",
		Port:        "5432",
		DBName:      "logtracker",
		User:        "raf",
		Password:    "qwq121",
		SslMode:     "disable",
	}

	servers["cute_ganymede"] = DBConfig{
		Name:        "cute_ganymede",
		DisplayName: "Cute Ganymede",
		DriverName:  "postgres",
		Host:        "147.45.226.19",
		Port:        "5432",
		DBName:      "logtracker",
		User:        "raf",
		Password:    "qwq121",
		SslMode:     "disable",
	}

	servers["test_server"] = DBConfig{
		Name:        "test_server",
		DisplayName: "Test Server",
		DriverName:  "postgres",
		Host:        "127.0.0.1",
		Port:        "5432",
		DBName:      "logtracker",
		User:        "raf",
		Password:    "qwq121",
		SslMode:     "disable",
	}

	for key, value := range servers {
		fmt.Println(key, ":", value)
	}

	file, err := os.Create("db-config.json")
	if err != nil {
		fmt.Println("error creating file")
		return
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(servers, "", "    ")
	if err != nil {
		fmt.Println("error marshaling json")
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("error writing to file")
		return
	}
}
