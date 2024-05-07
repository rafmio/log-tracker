package main

import (
	"encoding/json"
	"log"
	"os"
)

type ConnectDBConfig struct {
	DriverName string `json:"driverName"`
	User       string `json:"user"`
	Dbname     string `json:"dbname"`
	TableName  string `json:"tableName"`
	Password   string `json:"password"`
	Sslmode    string `json:"sslmode"`
}

func main() {
	var CDBc ConnectDBConfig = ConnectDBConfig{
		DriverName: "postgres",
		User:       "raf",
		Dbname:     "logtracker",
		TableName:  "lg_tab_1",
		Password:   "qwq121",
		Sslmode:    "disable",
	}

	file, err := os.Create("databaseConfig.json")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	// jsonData, err := json.Marshal(CDBc)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	jsonData, err := json.MarshalIndent(CDBc, "", "    ")
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("the file has been created")
}
