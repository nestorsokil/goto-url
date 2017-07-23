package config

import (
	"os"
	"encoding/json"
	"log"
	"path/filepath"
)

type Config struct {
	LogDir string `json:"log_dir"`
	Port string `json:"port"`
	MongoUrl string `json:"mongo_url"`
	Database string `json:"database"`
	KeyLength int `json:"key_length"`
	DevMode bool `json:"dev_mode"`
	ApplicationUrl string `json:"application_url"`
	ExpirationTimeHours int `json:"expiration_time_hours"`
	ClearTimeMinutes int `json:"clear_time_minutes"`
}

var Settings Config = loadConfig()

func loadConfig() Config {
	var conf Config
	file, err := os.Open("config/conf.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return conf
}

func GetGlobalLogFile() (*os.File) {
	logDir := filepath.Dir(Settings.LogDir)
	if logDir != "" {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.OpenFile(
		Settings.LogDir + string(os.PathSeparator)+
			"goto.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	return logFile
}

func GetRequestLogFile() (*os.File){
	logFile, err := os.OpenFile(
		Settings.LogDir + string(os.PathSeparator) +
			"request.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}
