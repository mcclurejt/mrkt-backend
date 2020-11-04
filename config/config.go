package config

import (
	"log"
	"os"
)

type ApiConfig struct {
	GlassNodeAPIKey string
	IEXCloudAPIKey  string
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
			GlassNodeAPIKey: getEnv("GLASSNODE_API_KEY", ""),
			IEXCloudAPIKey:  getEnv("IEX_CLOUD_API_KEY"),
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
