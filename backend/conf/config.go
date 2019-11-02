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
	EnvTraceDbEnabled   = "ENV_TRACE_DB_ENABLED"
	EnvRedisUrl         = "ENV_REDIS_URL"
	EnvRedisPass        = "ENV_REDIS_PASS"
)

const (
	InMemory = "inMemory"
	Redis    = "redis"
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

type TestConfig struct {
	Data map[string]interface{}
}

func (t *TestConfig) GetString(key string) string {
	return t.Data[key].(string)
}

func (t *TestConfig) GetInt(key string) int {
	return t.Data[key].(int)
}
