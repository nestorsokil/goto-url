package conf

import (
	"os"
	"strconv"
)

const (
	EnvPort             = "ENV_PORT"
	EnvKeyLen           = "ENV_KEY_LEN"
	EnvExpirationMillis = "ENV_EXPIRATION_MILLIS"
	EnvStorage          = "ENV_STORAGE"
	EnvRedisUrl         = "ENV_REDIS_URL"
	EnvRedisPass        = "ENV_REDIS_PASS"
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
