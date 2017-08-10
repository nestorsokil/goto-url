package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	MONGO     = "mongo"
	IN_MEMORY = "inMemory"
)

var configDirectory string

type ApplicationConfig struct {
	LogDir              string `json:"log_dir"`
	Port                string `json:"port"`
	KeyLength           int    `json:"key_length"`
	DevMode             bool   `json:"dev_mode"`
	Database            string `json:"database"`
	ApplicationUrl      string `json:"application_url"`
	ExpirationTimeHours int64  `json:"expiration_time_hours"`
	ClearTimeSeconds    int64  `json:"clear_time_seconds"`
}

type MongoConfig struct {
	MongoUrls     []string `json:"mongo_urls"`
	MongoUser     string   `json:"mongo_user"`
	MongoPassword string   `json:"mongo_password"`
	DatabaseName  string   `json:"database_name"`
	EnableTLS     bool     `json:"enable_tls"`
}

func LoadConfig() ApplicationConfig {
	configDirectory = os.Getenv("GO_TO_URL_CONFIG")
	if configDirectory == "" {
		configDirectory = "config/"
	}
	var conf ApplicationConfig
	file, err := os.Open(configDirectory + "conf.json")
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

func LoadMongoConfig() MongoConfig {
	var conf MongoConfig
	file, err := os.Open(configDirectory + "mongo_conf.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return conf
}

func (conf *ApplicationConfig) GetGlobalLogFile() *os.File {
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

func (conf *ApplicationConfig) GetRequestLogFile() *os.File {
	logFile, err := os.OpenFile(conf.LogDir+string(os.PathSeparator)+
		"request.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}
