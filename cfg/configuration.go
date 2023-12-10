package cfg

import "os"

type DbConnectionConfiguration interface {
	DbConnectionUrl() string
}

type JwtConfiguration interface {
	JwtSecret() string
}

type RedisConfiguration interface {
	RedisConnectionUrl() string
}

type Configuration interface {
	DbConnectionConfiguration
	JwtConfiguration
	RedisConfiguration
}

type ConfigurationKeys struct {
}

func (c ConfigurationKeys) DbConnectionUrl() string {
	return os.Getenv("DB_URL")
}

func (c ConfigurationKeys) JwtSecret() string {
	return os.Getenv("DB_URL")
}

func (c ConfigurationKeys) RedisConnectionUrl() string {
	return os.Getenv("REDIS_URL")
}
