package main

import (
	"os"
	"encoding/json"
	"log"
	"path/filepath"
)

type configuration struct {
	LogDir string `json:"log_dir"`
	Port string `json:"port"`
	MongoUrl string `json:"mongo_url"`
	Database string `json:"database"`
}

var conf configuration = loadConfig()

func loadConfig() configuration {
	var conf configuration
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

func setupGlobalLog() {
	logDir := filepath.Dir(conf.LogDir)
	if logDir != "" {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.OpenFile(
		conf.LogDir + string(os.PathSeparator) + "goto.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
}

func getRequestLogFile() (*os.File){
	logFile, err := os.OpenFile(
		conf.LogDir + string(os.PathSeparator) +
			"request.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}
