package config

import (
	"os"
	"strconv"
)

var (
	Port, _    = strconv.Atoi(getEnv("APP_PORT", "8080"))
	AppVersion = getEnv("APP_VERSION", "MISSING")
	AppName    = getEnv("APP_NAME", "film-list")
)

func getEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}
