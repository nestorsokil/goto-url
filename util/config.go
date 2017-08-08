package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	LogDir              string   `json:"log_dir"`
	Port                string   `json:"port"`
	EnableTLS           bool     `json:"enable_tls"`
	MongoUrls           []string `json:"mongo_urls"`
	MongoUser           string   `json:"mongo_user"`
	MongoPassword       string   `json:"mongo_password"`
	Database            string   `json:"database"`
	KeyLength           int      `json:"key_length"`
	DevMode             bool     `json:"dev_mode"`
	ApplicationUrl      string   `json:"application_url"`
	ExpirationTimeHours int64    `json:"expiration_time_hours"`
	ClearTimeSeconds    int64    `json:"clear_time_seconds"`
}

func LoadConfig() Configuration {
	confPath := os.Getenv("GO_TO_URL_CONFIG")
	if confPath == "" {
		confPath = "config/conf.json"
	}
	var conf Configuration
	file, err := os.Open(confPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return conf
}

func (conf *Configuration) GetGlobalLogFile() *os.File {
	logDir := filepath.Dir(conf.LogDir)
	if logDir != "" {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.OpenFile(conf.LogDir+string(os.PathSeparator)+
		"goto.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}

func (conf *Configuration) GetRequestLogFile() *os.File {
	logFile, err := os.OpenFile(conf.LogDir+string(os.PathSeparator)+
		"request.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}
