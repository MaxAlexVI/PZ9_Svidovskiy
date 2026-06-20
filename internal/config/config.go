package config

import "time"

type Config struct {
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	TaskCacheTTL    time.Duration
	TaskCacheJitter time.Duration
	ServerAddr      string
}

func New() Config {
	return Config{
		RedisAddr:       "localhost:6379",
		RedisPassword:   "",
		RedisDB:         0,
		TaskCacheTTL:    30 * time.Second,
		TaskCacheJitter: 10 * time.Second,
		ServerAddr:      ":8082",
	}
}
