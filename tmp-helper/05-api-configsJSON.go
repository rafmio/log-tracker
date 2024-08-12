package main

type Server struct {
	Name      string `json:"Name"`
	IP        string `json:"IP"`
	Port      string `json:"Port"`
	DBName    string `json:"DBName"`
	TableName string `json:"TableName"`
	Password  string `json:"Password"`
}
