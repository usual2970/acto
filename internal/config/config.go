package config

import (
	"os"
)

type Config struct {
	MySQLDSN  string
	RedisAddr string
	HTTPAddr  string
}

func Load() Config {
	cfg := Config{
		MySQLDSN:  getenv("MYSQL_DSN", "acto:acto@tcp(127.0.0.1:3306)/acto?parseTime=true&charset=utf8mb4&loc=Local"),
		RedisAddr: getenv("REDIS_ADDR", "127.0.0.1:6379"),
		HTTPAddr:  getenv("HTTP_ADDR", ":8080"),
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
