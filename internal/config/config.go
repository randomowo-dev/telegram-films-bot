package config

import (
	"os"
)

var (
	AppVersion = getEnv("APP_VERSION", "MISSING")
	AppName    = getEnv("APP_NAME", "film-list")
)

func getEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}
