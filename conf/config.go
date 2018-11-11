package conf

import (
	"os"
	"strconv"
)

const (
	ENV_LOG_DIR           = "ENV_LOG_DIR"
	ENV_PORT              = "ENV_PORT"
	ENV_KEY_LEN           = "ENV_KEY_LEN"
	ENV_EXPIRATION_MILLIS = "ENV_EXPIRATION_MILLIS"
	ENV_STORAGE           = "ENV_STORAGE"
	ENV_REDIS_URL         = "ENV_REDIS_URL"
)

const (
	IN_MEMORY = "inMemory"
	REDIS     = "redis"
)

type Config interface {
	GetString(string) string
	GetInt(string) int
}

type EnvConfig struct{}

func (e *EnvConfig) GetString(key string) string {
	return os.Getenv(key)
}

func (e *EnvConfig) GetInt(key string) int {
	i, _ := strconv.Atoi(os.Getenv(key))
	return i
}
