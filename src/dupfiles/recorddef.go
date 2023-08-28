package main

type Syslog struct {
	RemoteHost string `json:"remotehost"`
	Port       int    `json:"port"`
	Tag        string `json:"tag"`
	Protocol   string `json:"protocol"`
}

type DB struct {
	DBName   string `json:"dbname"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Config struct {
	Database DB
	Syslog   Syslog
}
