package config

import (
	"log"
	"os"
)

type ApiConfig struct {
	AlphavantageAPIKey string
	MarketStackAPIKey  string
}

type DbConfig struct {
	Datasource string
}

type Config struct {
	Api ApiConfig
	Db  DbConfig
}

func New() *Config {
	return &Config{
		Api: ApiConfig{
			AlphavantageAPIKey: getEnv("ALPHAVANTAGE_API_KEY", ""),
			MarketStackAPIKey:  getEnv("MARKETSTACK_API_KEY", ""),
		},
		Db: DbConfig{
			Datasource: getEnv("DB_DATASOURCE", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("[Warning]: There is no env variable for %s.  Defaulting to \"\"", key)
	return defaultVal
}
