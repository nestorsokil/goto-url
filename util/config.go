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
	REDIS     = "redis"
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

type RedisConfig struct {
	RedisUrl string `json:"redis_url"`
}

func LoadConfig() ApplicationConfig {
	configDirectory = os.Getenv("GO_TO_URL_CONFIG")
	if configDirectory == "" {
		configDirectory = "config/"
	}
	var conf ApplicationConfig
	configPath := configDirectory + "conf.json"
	parseConfig(configPath, &conf)
	return conf
}

func LoadMongoConfig() MongoConfig {
	var conf MongoConfig
	configPath := configDirectory + "mongo_conf.json"
	parseConfig(configPath, &conf)
	return conf
}

func LoadRedisConfig() RedisConfig {
	var conf RedisConfig
	configPath := configDirectory + "redis_conf.json"
	parseConfig(configPath, &conf)
	return conf
}

func parseConfig(fromFile string, toStruct interface{}) {
	file, err := os.Open(fromFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(toStruct)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

func (conf *ApplicationConfig) GetRequestLogFile() *os.File {
	logDir := filepath.Dir(conf.LogDir)
	if logDir != "" {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.OpenFile(conf.LogDir+string(os.PathSeparator)+
		"request.log", os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}
