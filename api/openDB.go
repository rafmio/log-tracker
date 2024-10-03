package main

import (
	"encoding/json"
	"log"
	"os"
)

type ConnectDBConfig struct {
	DriverName string `json:"DriverName"`
	Name       string `json:"Name"`
	Host       string `json:"Host"`
	Port       string `json:"Port"`
	DBName     string `json:"DBName"`
	User       string `json:"User"`
	TableName  string `json:"TableName"`
	Password   string `json:"Password"`
	SslMode    string `json:"SslMode"`
}

func setDbConfigFilePath() string {
	return "db-config.json"
}

func readConfig(dbConfigFilePath string) (map[string]ConnectDBConfig, error) {
	// reading file with configuration for DB connection
	file, err := os.ReadFile(dbConfigFilePath)
	if err != nil {
		log.Println("Opening config file:", err)
	}

	// unmarshalling JSON data to struct
	dbConfigs := make(map[string]ConnectDBConfig)

	err = json.Unmarshal(file, &dbConfigs)
	if err != nil {
		log.Println("Unmarshalling JSON:", err)
	}

	return dbConfigs, err
}
