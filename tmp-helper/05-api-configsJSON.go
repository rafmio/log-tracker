package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Server struct {
	Name      string `json:"Name"`
	Host      string `json:"Host"`
	Port      string `json:"Port"`
	DBName    string `json:"DBName"`
	User      string `json:"User"`
	TableName string `json:"TableName"`
	Password  string `json:"Password"`
	SslMode   string `json:"SslMode"`
}

func main() {
	servers := make(map[string]Server, 3)

	servers["BlackOxygenium"] = Server{
		Name:      "BlackOxygenium",
		Host:      "194.58.102.129",
		Port:      "5432",
		DBName:    "logtracker",
		User:      "raf",
		TableName: "lg_tab",
		Password:  "qwq121",
		SslMode:   "disable",
	}

	servers["CuteGanymede"] = Server{
		Name:      "CuteGanymede",
		Host:      "147.45.226.19",
		Port:      "5432",
		DBName:    "logtracker",
		User:      "raf",
		TableName: "lg_tag",
		Password:  "qwq121",
		SslMode:   "disable",
	}

	servers["TestServer"] = Server{
		Name:      "TestServer",
		Host:      "127.0.0.1",
		Port:      "5432",
		DBName:    "logtracker",
		User:      "raf",
		TableName: "lg_tag",
		Password:  "qwq121",
		SslMode:   "disable",
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
